package main

import (
	"log"
	"order_service/composer"
	"order_service/config"
	"order_service/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetUpRoutes(router fiber.Router, cfg *config.Config, db *pgxpool.Pool) {
	// create businesses
	authBiz := composer.ComposeAuthBusiness(cfg, db)
	userBiz := composer.ComposeUserBusiness(cfg, db)

	// create services
	authAPIService := composer.ComposeAuthAPIService(authBiz)
	userAPIService := composer.ComposeUserAPIService(userBiz)

	// create middlewares
	authMiddleware := middleware.RequireAuth(authBiz)

	// prepare routes
	// /auth
	authRouter := router.Group("/auth")
	{

		authRouter.Post("/register", authAPIService.Register)
		authRouter.Post("/login", authAPIService.Login)
	}

	// /users
	userRouter := router.Group("/users", authMiddleware)
	{
		userRouter.Get("/profile", userAPIService.GetUserProfile)
	}
}

func main() {
	cfg := config.NewConfig()
	db := config.ConnectToPostgres(cfg)
	defer db.Close()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())

	SetUpRoutes(app.Group("/v1"), cfg, db)

	log.Fatalln(app.Listen("0.0.0.0:8080"))
}
