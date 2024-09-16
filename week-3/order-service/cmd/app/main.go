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
	authUc := composer.ComposeAuthUsecase(cfg, pg, rd)
	userUc := composer.ComposeUserUsecase(pg)
	productUc := composer.ComposeProductUsecase(pg)
	orderUc := composer.ComposeOrderUsecase(pg)

	// create services
	authAPIService := composer.ComposeAuthAPIService(authUc)
	userAPIService := composer.ComposeUserAPIService(userUc)
	productAPIService := composer.ComposeProductAPIService(productUc)
	orderAPIService := composer.ComposeOrderAPIService(orderUc)

	// create middlewares
	authMiddleware := middleware.RequireAuth(authUc)

	// prepare routes
	// /auth
	authRouter := router.Group("/auth")
	{

		authRouter.Post("/register", authAPIService.Register)
		authRouter.Post("/login", authAPIService.Login)
		authRouter.Post("/refresh", authAPIService.Refresh)
		authRouter.Post("/sign-out", authMiddleware, authAPIService.SignOut)
		authRouter.Post("/sign-out-all", authMiddleware, authAPIService.SignOutAll)
	}

	// /users
	userRouter := router.Group("/users", authMiddleware)
	{
		userRouter.Get("/", userAPIService.GetUsers)
		userRouter.Get("/:userID", userAPIService.GetUser)
		userRouter.Get("/profile", userAPIService.GetUserProfile)
		userRouter.Post("/balance", userAPIService.AddUserBalance)
	}

	// /products
	productRouter := router.Group("/products")
	{
		productRouter.Get("/", productAPIService.GetProducts)
		productRouter.Get("/:productID", productAPIService.GetProduct)
		productRouter.Post("/", authMiddleware, productAPIService.CreateProduct)
		productRouter.Put("/:productID", authMiddleware, productAPIService.UpdateProduct)
		productRouter.Delete("/:productID", authMiddleware, productAPIService.DeleteProduct)
	}

	// /orders
	orderRouter := router.Group("/orders", authMiddleware)
	{
		orderRouter.Get("/", orderAPIService.GetOrders)
		orderRouter.Post("/", orderAPIService.CreateOrder)
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
