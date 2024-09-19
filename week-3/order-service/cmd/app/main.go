package main

import (
	"log"
	"order_service/composer"
	"order_service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.NewConfig()
	pg := config.ConnectToPostgres(cfg)
	rd := config.ConnectToRedis(cfg)
	defer pg.Close()
	defer rd.Close()

	app := fiber.New()

	app.Use(recover.New())
	// app.Use(limiter.New(limiter.Config{
	// 	Max:        256,
	// 	Expiration: 5 * time.Second,
	// }))

	composer.SetUpRoutes(app.Group("/v1"), cfg, pg, rd)

	log.Fatalln(app.Listen("0.0.0.0:8080"))
}
