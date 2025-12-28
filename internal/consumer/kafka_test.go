package consumer

import (
	"context"
	"testing"

	"rating-service/internal/grpc"
	"rating-service/internal/service"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	gamepb "github.com/viktoralyoshin/playhub-proto/gen/go/games"
	socialpb "github.com/viktoralyoshin/playhub-proto/gen/go/social"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockKafkaReader struct {
	mock.Mock
}

func (m *MockKafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockKafkaReader) Close() error {
	return m.Called().Error(0)
}

func (m *MockKafkaReader) Config() kafka.ReaderConfig {
	return kafka.ReaderConfig{Topic: "test-topic", GroupID: "test-group"}
}

type MockSocialClient struct {
	mock.Mock
	socialpb.SocialServiceClient
}

func (m *MockSocialClient) GetGameReviews(ctx context.Context, in *socialpb.GetGameReviewsRequest, opts ...ggrpc.CallOption) (*socialpb.GetGameReviewsResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*socialpb.GetGameReviewsResponse), args.Error(1)
}

type MockGamesClient struct {
	mock.Mock
	gamepb.GameServiceClient
}

func TestNewRatingConsumer(t *testing.T) {
	broker := "localhost:9092"
	topic := "test-topic"
	svc := service.NewRatingService()

	consumer := NewRatingConsumer(broker, topic, svc)

	assert.NotNil(t, consumer)
	assert.NotNil(t, consumer.reader)
	assert.Equal(t, svc, consumer.service)

	assert.Equal(t, topic, consumer.reader.Config().Topic)
	assert.Contains(t, consumer.reader.Config().Brokers, broker)

	err := consumer.reader.Close()
	assert.NoError(t, err)
}

func (m *MockGamesClient) SetRating(ctx context.Context, in *gamepb.RatingRequest, opts ...ggrpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func TestRatingConsumer_Run(t *testing.T) {
	gameID := uuid.New().String()

	oldSocial := grpc.SocialClient
	oldGames := grpc.GamesClient
	mockSocial := new(MockSocialClient)
	mockGames := new(MockGamesClient)
	grpc.SocialClient = mockSocial
	grpc.GamesClient = mockGames
	defer func() {
		grpc.SocialClient = oldSocial
		grpc.GamesClient = oldGames
	}()

	svc := service.NewRatingService()

	t.Run("full success cycle", func(t *testing.T) {
		mockReader := new(MockKafkaReader)
		c := &RatingConsumer{
			reader:  mockReader,
			service: svc,
		}

		ctx, cancel := context.WithCancel(context.Background())

		mockReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
			Value:  []byte(`{"game_id":"` + gameID + `"}`),
			Offset: 123,
		}, nil).Once()

		mockSocial.On("GetGameReviews", mock.Anything, mock.Anything).Return(&socialpb.GetGameReviewsResponse{
			Reviews: []*socialpb.Review{{Rating: 100}, {Rating: 80}},
		}, nil).Once()

		mockGames.On("SetRating", mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil).Once()

		mockReader.On("ReadMessage", mock.Anything).Return(kafka.Message{}, context.Canceled).Run(func(args mock.Arguments) {
			cancel()
		}).Once()

		mockReader.On("Close").Return(nil).Once()

		c.Run(ctx)

		mockReader.AssertExpectations(t)
		mockSocial.AssertExpectations(t)
		mockGames.AssertExpectations(t)
	})

	t.Run("unmarshal error and exit", func(t *testing.T) {
		mockReader := new(MockKafkaReader)
		c := &RatingConsumer{reader: mockReader, service: svc}
		ctx, cancel := context.WithCancel(context.Background())

		mockReader.On("ReadMessage", mock.Anything).Return(kafka.Message{
			Value: []byte(`{bad json}`),
		}, nil).Once()

		mockReader.On("ReadMessage", mock.Anything).Return(kafka.Message{}, context.Canceled).Run(func(args mock.Arguments) {
			cancel()
		}).Once()

		mockReader.On("Close").Return(nil).Once()

		c.Run(ctx)
		mockReader.AssertExpectations(t)
	})
}
