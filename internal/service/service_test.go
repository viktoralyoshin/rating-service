package service

import (
	"context"
	"errors"
	"testing"

	"rating-service/internal/grpc"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viktoralyoshin/playhub-proto/gen/go/games"
	socialpb "github.com/viktoralyoshin/playhub-proto/gen/go/social"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
	games.GameServiceClient
}

func (m *MockGamesClient) SetRating(ctx context.Context, in *games.RatingRequest, opts ...ggrpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func TestRatingService_CalculateRating(t *testing.T) {
	oldSocialClient := grpc.SocialClient
	mockSocial := new(MockSocialClient)
	grpc.SocialClient = mockSocial
	defer func() { grpc.SocialClient = oldSocialClient }()

	svc := NewRatingService()
	gameID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockSocial.On("GetGameReviews", mock.Anything, mock.Anything).Return(&socialpb.GetGameReviewsResponse{
			Reviews: []*socialpb.Review{
				{Rating: 10},
				{Rating: 20},
			},
		}, nil).Once()

		res, err := svc.CalculateRating(context.Background(), gameID)
		assert.NoError(t, err)
		assert.Equal(t, 15, res)
	})

	t.Run("grpc error", func(t *testing.T) {
		mockSocial.On("GetGameReviews", mock.Anything, mock.Anything).
			Return(nil, errors.New("connection lost")).Once()

		res, err := svc.CalculateRating(context.Background(), gameID)
		assert.Error(t, err)
		assert.Equal(t, 0, res)
	})
}

func TestRatingService_SendRating(t *testing.T) {
	oldGamesClient := grpc.GamesClient
	mockGames := new(MockGamesClient)
	grpc.GamesClient = mockGames
	defer func() { grpc.GamesClient = oldGamesClient }()

	svc := NewRatingService()
	gameID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockGames.On("SetRating", mock.Anything, &games.RatingRequest{
			GameId: gameID,
			Rating: 90,
		}).Return(&emptypb.Empty{}, nil).Once()

		err := svc.SendRating(context.Background(), gameID, 90)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockGames.On("SetRating", mock.Anything, mock.Anything).
			Return(nil, errors.New("fail")).Once()

		err := svc.SendRating(context.Background(), gameID, 90)
		assert.Error(t, err)
	})
}
