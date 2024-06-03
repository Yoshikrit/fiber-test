package router

import (
	"github.com/Yoshikrit/fiber-test/middleware"
	"github.com/Yoshikrit/fiber-test/handler"
	"github.com/Yoshikrit/fiber-test/service"
	"github.com/Yoshikrit/fiber-test/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/gofiber/swagger"
	_ "github.com/Yoshikrit/fiber-test/docs"
)

func NewRouter(db *gorm.DB) *fiber.App {
	router := fiber.New()

	router.Get("/swagger/*", swagger.HandlerDefault)
	router.Get("/metrics", middleware.Metrics())

	//healthcheck
	healthCheckHandler := handler.NewHealthCheckHandler()

	router.Get("/healthcheck", healthCheckHandler.HealthCheck)

	//producttype
	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)

	productTypeRouter := router.Group("/producttype")
	
	productTypeRouter.Post("/", prodTypeHandler.Create)
	productTypeRouter.Get("/", prodTypeHandler.FindAll)
	productTypeRouter.Get("/count", prodTypeHandler.Count)

	productTypeRouter.Route("/:id", func(router fiber.Router) {
		router.Get("/", prodTypeHandler.FindByID)
		router.Put("/", prodTypeHandler.Update)
		router.Delete("/", prodTypeHandler.Delete)
	})


	return router
}