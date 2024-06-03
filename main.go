package main

import (
	"github.com/Yoshikrit/fiber-test/router"
	"github.com/Yoshikrit/fiber-test/config"
	"github.com/Yoshikrit/fiber-test/middleware"
	"github.com/Yoshikrit/fiber-test/model"
	
	"github.com/gofiber/fiber/v2"
)

// @title ProductType API for For Fiber-Test
// @description API ProductType management Server by Fiber- Teletubbie's ProductType API.
// @version 1.0

// @contact.name   Walter White
// @contact.url    https://twitter.com/example
// @contact.email  example@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /

// @schemes http https

// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
func main() {
	app := fiber.New()
	
	//config
	loadConfig, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	config.InitTimeZone()

	//Database
	db := config.ConnectionDB(&loadConfig)
	db.AutoMigrate(&model.ProductTypeEntity{})

	//Routes
	routes := router.NewRouter(db)

	//middleware
	app.Use(
		middleware.Cors(), 
		middleware.Recover(), 
		middleware.Logger(), 
		middleware.Health(),
		middleware.Limiter(),
	)
	app.Mount("/", routes)

	app.Listen(":" +  loadConfig.ServerPort)
}