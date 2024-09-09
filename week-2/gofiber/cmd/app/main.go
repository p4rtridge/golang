package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/partridge1307/gofiber/config"
	"github.com/partridge1307/gofiber/controllers"
	"github.com/partridge1307/gofiber/infrastructure"
	"github.com/partridge1307/gofiber/usecases"
)

func main() {
	// Init config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("[Error]: Config error: %s", err.Error())
	}

	// Create Fiber's application
	app := fiber.New(fiber.Config{
		AppName:      cfg.Name,
		ServerHeader: cfg.Header,
	})

	// Add middlewares
	app.Use(logger.New())

	// Create repos
	postgresRepo, err := infrastructure.NewPostgresRepo(cfg)
	if err != nil {
		log.Fatalf("[Error]: Init repo error: %s", err.Error())
	}

	// Create use cases
	authUsecase := usecases.NewAuthUsecase(postgresRepo)
	userUsecase := usecases.NewUserUsecase(postgresRepo)

	// Prepare endpoints
	controllers.NewAuthController(app.Group("/api/v1/auth"), authUsecase)
	controllers.NewUserController(app.Group("/api/v1/users"), userUsecase)

	// Listen to port 8080
	log.Fatal(app.Listen("0.0.0.0:8080"))
}
