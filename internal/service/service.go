package service

import (
	"context"
	"rating-service/internal/grpc"
	"rating-service/internal/utils"

	"github.com/viktoralyoshin/playhub-proto/gen/go/games"
	socialpb "github.com/viktoralyoshin/playhub-proto/gen/go/social"
)

type RatingService struct{}

func NewRatingService() *RatingService {
	return &RatingService{}
}

func (s *RatingService) CalculateRating(ctx context.Context, gameId string) (int, error) {
	resp, err := grpc.SocialClient.GetGameReviews(ctx, &socialpb.GetGameReviewsRequest{
		GameId: gameId,
		Limit:  0,
		Offset: 0,
	})
	if err != nil {
		return 0, err
	}

	ratings := make([]int, 0, len(resp.Reviews))

	for _, review := range resp.Reviews {
		ratings = append(ratings, int(review.Rating))
	}

	return utils.CalculateAverage(ratings), nil

}

func (s *RatingService) SendRating(ctx context.Context, gameId string, rating int) error {

	_, err := grpc.GamesClient.SetRating(ctx, &games.RatingRequest{
		GameId: gameId,
		Rating: uint32(rating),
	})
	if err != nil {
		return err
	}

	return nil
}
