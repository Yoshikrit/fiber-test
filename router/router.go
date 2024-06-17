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

func NewRouter(router *fiber.App, db *gorm.DB) *fiber.App {
	router.Get("/swagger/*", swagger.HandlerDefault)
	router.Get("/metrics", middleware.Metrics())

	//healthcheck
	healthCheckHandler := handler.NewHealthCheckHandler()

	router.Get("/healthcheck", healthCheckHandler.HealthCheck)

	//auths
	userRepository := repository.NewUserRepositoryImpl(db)
	roleRepository := repository.NewRoleRepositoryImpl(db)
	oauthRepository := repository.NewOauthRepositoryImpl(db)
	authService := service.NewAuthServiceImpl(userRepository, roleRepository, oauthRepository)
	authHandler := handler.NewAuthHandler(authService)

	authRouter := router.Group("/auths")

	authRouter.Post("/", authHandler.Register)
	authRouter.Post("/login", authHandler.Login)
	authRouter.Post("/reflesh", authHandler.Reflesh)
	authRouter.Route("logout/:id", func(router fiber.Router) {
		router.Delete("/", authHandler.Logout)
	})

	//create jwt middleware
	jwtMiddleware := middleware.NewJWTMiddleware(userRepository, oauthRepository, roleRepository)

	//producttypes
	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)

    productTypeRouter := router.Group("/producttypes")
	productTypeRouter.Use(jwtMiddleware)
	
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