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
	"github.com/redis/go-redis/v9"
)

func SetUpRoutes(router fiber.Router, cfg *config.Config, pg *pgxpool.Pool, rd *redis.Client) {
	// create businesses
	authBiz := composer.ComposeAuthBusiness(cfg, pg, rd)
	userBiz := composer.ComposeUserBusiness(cfg, pg)

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
		authRouter.Post("/refresh", authAPIService.Refresh)
		authRouter.Post("/refresh", authAPIService.Refresh)
		authRouter.Post("/sign-out", authMiddleware, authAPIService.SignOut)
		authRouter.Post("/sign-out-all", authMiddleware, authAPIService.SignOutAll)
	}

	// /users
	userRouter := router.Group("/users", authMiddleware)
	{
		userRouter.Get("/profile", userAPIService.GetUserProfile)
	}
}

func main() {
	cfg := config.NewConfig()
	pg := config.ConnectToPostgres(cfg)
	rd := config.ConnectToRedis(cfg)
	defer pg.Close()
	defer rd.Close()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())

	SetUpRoutes(app.Group("/v1"), cfg, pg, rd)

	log.Fatalln(app.Listen("0.0.0.0:8080"))
}
