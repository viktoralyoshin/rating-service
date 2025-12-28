package consumer

import (
	"context"
	"encoding/json"
	"rating-service/internal/service"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
	Config() kafka.ReaderConfig
}

type ReviewEvent struct {
	GameID string `json:"game_id"`
}

type RatingConsumer struct {
	reader  KafkaReader
	service *service.RatingService
}

func NewRatingConsumer(broker string, topic string, service *service.RatingService) *RatingConsumer {
	return &RatingConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{broker},
			Topic:       topic,
			GroupID:     "rating-group",
			StartOffset: kafka.FirstOffset,
		}),
		service: service,
	}
}

func (c *RatingConsumer) Run(ctx context.Context) {
	log.Info().
		Str("topic", c.reader.Config().Topic).
		Str("group_id", c.reader.Config().GroupID).
		Msg("Rating Consumer started and waiting messages")

	defer func() {
		if err := c.reader.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close kafka reader")
		}
	}()

	for {
		message, err := c.reader.ReadMessage(ctx)
		log.Info().Msg("Kafka message recieved")
		if err != nil {
			if ctx.Err() != nil {
				log.Info().Msg("consumer context cancelled")

				return
			}
			log.Error().Err(err).Msg("failed to read message from kafka")
			continue
		}

		var event ReviewEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Error().
				Err(err).
				Str("payload", string(message.Value)).
				Msg("failed to unmarshal kafka message")

			continue
		}

		l := log.With().
			Str("game_id", event.GameID).
			Int64("offset", message.Offset).
			Logger()

		rating, err := c.service.CalculateRating(ctx, event.GameID)
		if err != nil {
			l.Error().Err(err).Msg("failed to calculate rating")

			continue
		}

		l.Debug().Int("calculated_rating", rating).Msg("rating calculated")

		if err := c.service.SendRating(ctx, event.GameID, rating); err != nil {
			l.Error().Err(err).Msg("failed to send rating update")

			continue
		}

		l.Info().
			Int("rating", rating).
			Msg("rating successfully updated")
	}
}
