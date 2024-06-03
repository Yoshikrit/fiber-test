package middleware

import (
	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	config := cors.Config{
		AllowOrigins: "http://localhost:8081, https://localhost:8081",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}

	return cors.New(config)
}