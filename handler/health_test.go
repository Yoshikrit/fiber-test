package handler_test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"testing"
	"io"

	"github.com/Yoshikrit/fiber-test/handler"
)

func TestHealthCheck(t *testing.T) {
	healthHandler := handler.NewHealthCheckHandler()

	app := fiber.New()
  	app.Get("/health", healthHandler.HealthCheck)
	t.Run("test case : pass", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":"Welcome to ProductType Server"}`
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, expectedBody, string(body))
	})
}