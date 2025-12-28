package router

import (
	"rating-service/internal/handler"
	"rating-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	rating := v1.Group("/games")

	ratingService := service.NewRatingService()
	ratingHandler := handler.NewRatingHandler(ratingService)

	rating.Get("/calculate", ratingHandler.Calculate)
}
