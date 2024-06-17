package helper_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/gofiber/fiber/v2"

	"github.com/Yoshikrit/fiber-test/helper"

	"net/http"
	"net/http/httptest"
	"testing"
	"io"
)


func TestParamsInt(t *testing.T) {
	app := fiber.New()

	app.Get("/int/:id", func(ctx *fiber.Ctx) error {
		id, err := helper.ParamsInt(ctx)
		if err != nil {
			return helper.HandleError(ctx, err)
		}
		return ctx.JSON(fiber.Map{"id": id})
	})

	t.Run("Valid Integer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/int/123", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		expectedBody := `{"id":123}`
		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, expectedBody, string(body))
	})

	t.Run("Invalid Non-Integer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/int/abc", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Invalid ID: abc is not integer"}`
		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, expectedBody, string(body))
	})
}
