package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Yoshikrit/fiber-test/helper/logger"
	"github.com/Yoshikrit/fiber-test/model"
)

type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// HealthCheckHandler godoc
// @Summary Health Check
// @Description Health check
// @id HealthCheckHandler
// @Tags healthcheck
// @Produce  json
// @response 200 {object} model.StringResponse "Welcome to ProductType Server"
// @Router /healthcheck [get]
func (h *HealthCheckHandler) HealthCheck(ctx *fiber.Ctx) error {
	logger.Info("Handler: ProductType service is running")
	webResponse := model.StringResponse{
		Code: 		200,
		Message: 	"Welcome to ProductType Server",
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}
