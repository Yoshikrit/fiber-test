package middleware

import (
	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func Health() fiber.Handler {
	return healthcheck.New()
}