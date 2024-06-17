package main

import (
	"github.com/Yoshikrit/fiber-test/router"
	"github.com/Yoshikrit/fiber-test/config"
	"github.com/Yoshikrit/fiber-test/middleware"
	// "github.com/Yoshikrit/fiber-test/model"
	
	"github.com/gofiber/fiber/v2"
	"github.com/goccy/go-json"
)

// @title ProductType API for Fiber-Test
// @description API ProductType management Server by Fiber-Teletubbie's ProductType API.
// @version 1.0

// @contact.name   Walter White
// @contact.url    https://twitter.com/example
// @contact.email  example@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /

// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Type 'Bearer' followed by a space and your JWT token."
func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	
	//config
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	configData, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	//Database
	db := config.ConnectionDB(&configData)
	// db.AutoMigrate(&model.ProductTypeEntity{}, &models.UserEntity{}, &models.RoleEntity{}, &models.OauthEntity{})

	//Routes
	router.NewRouter(app, db)

	//middleware
	app.Use(
		middleware.Cors(), 
		middleware.Recover(), 
		middleware.Logger(), 
		middleware.Health(),
		middleware.Limiter(),
		middleware.Helmet(),
	)

	app.Listen(":" +  configData.ServerPort)
}