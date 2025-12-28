package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRatingHandler_Calculate(t *testing.T) {
	app := fiber.New()
	h := NewRatingHandler(nil)

	app.Get("/calculate", h.Calculate)

	req := httptest.NewRequest("GET", "/calculate", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}
