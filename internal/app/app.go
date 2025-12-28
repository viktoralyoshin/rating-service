package app

import (
	"context"
	"rating-service/internal/config"
	"rating-service/internal/consumer"
	"rating-service/internal/grpc"
	"rating-service/internal/router"
	"rating-service/internal/service"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Start(cfg *config.Config) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info().Msg("Rating Service pre-start...")
	grpc.Connect(cfg)
	log.Info().Msg("gRPC connected...")

	log.Info().
		Msg("creating consumer...")

	consumer := consumer.NewRatingConsumer(
		cfg.KafkaAddr,
		"review_events",
		service.NewRatingService(),
	)

	log.Info().
		Msg("run consumer...")

	go consumer.Run(ctx)

	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	router.SetupRouter(app)

	port := ":" + cfg.Port

	log.Info().
		Str("port", port).
		Msg("Rating Service starting")

	if err := app.Listen(port); err != nil {
		log.Fatal().Err(err).Msg("Rating Service server failed")
	}

}
