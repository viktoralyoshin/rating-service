package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	app := fiber.New()
	SetupRouter(app)

	req := httptest.NewRequest("GET", "/api/v1/games/calculate", nil)
	resp, _ := app.Test(req)

	assert.NotEqual(t, 404, resp.StatusCode)
}
