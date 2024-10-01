package composer

import (
	"order_service/config"
	"order_service/middleware"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetUpRoutes(router fiber.Router, cfg *config.Config, pg *pgxpool.Pool, rd *redis.Client, s3Client *s3.Client) {
	// create businesses
	authUc := ComposeAuthUsecase(cfg, pg, rd)
	userUc := ComposeUserUsecase(pg)
	productUc := ComposeProductUsecase(pg, s3Client)
	orderUc := ComposeOrderUsecase(pg)

	// create services
	authAPIService := ComposeAuthAPIService(authUc)
	userAPIService := ComposeUserAPIService(userUc)
	productAPIService := ComposeProductAPIService(productUc)
	orderAPIService := ComposeOrderAPIService(orderUc)

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
		userRouter.Get("/profile", userAPIService.GetUserProfile)
		userRouter.Get("/:userID", userAPIService.GetUser)
		userRouter.Post("/balance", userAPIService.AddUserBalance)
	}

	// /products
	productRouter := router.Group("/products")
	{
		productRouter.Get("/", productAPIService.GetProducts)
		productRouter.Get("/search/", productAPIService.SearchProducts)
		productRouter.Get("/:productID", productAPIService.GetProduct)
		productRouter.Post("/", authMiddleware, productAPIService.CreateProduct)
		productRouter.Put("/:productID", authMiddleware, productAPIService.UpdateProduct)
		productRouter.Delete("/:productID", authMiddleware, productAPIService.DeleteProduct)
	}

	// /orders
	orderRouter := router.Group("/orders", authMiddleware)
	{
		orderRouter.Get("/", orderAPIService.GetOrders)
		orderRouter.Get("/top-by-price", orderAPIService.GetTopFiveOrdersByPrice)
		orderRouter.Get("/orders-by-month", orderAPIService.GetNumOfOrdersByMonth)
		orderRouter.Get("/:orderID", orderAPIService.GetOrder)
		orderRouter.Post("/", orderAPIService.CreateOrder)
		orderRouter.Post("/summarize", orderAPIService.GetOrdersSummarize)
	}
}
