package main

import (
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/partridge1307/gofiber/api/handler"
	"github.com/partridge1307/gofiber/infras/postgres"
	"github.com/partridge1307/gofiber/pkg/validate"
	"github.com/partridge1307/gofiber/usecase/auth"
	"github.com/partridge1307/gofiber/usecase/user"
)

func main() {
	pgPool, err := postgres.ConnectToPostgres(os.Getenv("PG_URL"))
	if err != nil {
		log.Fatalln("Failed to connect to postgres")
	}

	// Create support packages
	validator := validate.New()

	// Create repositories
	postgresRepo := postgres.NewPostgresRepo(pgPool)

	// Create usecases
	authService := auth.NewService(validator, postgresRepo)
	userService := user.NewService(postgresRepo)

	// Create fiber's app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(logger.New())

	// Prepare routes
	handler.NewAuthHandler(app.Group("/api/v1/auth"), authService)
	handler.NewUserHandler(app.Group("/api/v1/users"), userService)

	log.Fatalln(app.Listen("0.0.0.0:8080"))
}
