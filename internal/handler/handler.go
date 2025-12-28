package handler

import (
	"rating-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

type RatingHandler struct {
	service *service.RatingService
}

func NewRatingHandler(service *service.RatingService) *RatingHandler {
	return &RatingHandler{
		service: service,
	}
}

func (h *RatingHandler) Calculate(c *fiber.Ctx) error {
	return nil
}
