package main

import (
	_ "api_document/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title API document example
// @version 1
// @description This is a example server for API document purpose
// @termsOfService  http://swagger.io/terms/

// @contact.name partridge
// @contact.url https://iloveyour.dad
// @contact.email anhduc130703@gmail.com

// @license.name MIT
// @license.url https://opensource.org/license/mit

// @host localhost:8080
// @BasePath /v1
func main() {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", GetHandle)

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendStatus(201)
	})
	app.Put("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Delete("/", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	app.Listen(":8080")
}

// GetOrderByCode Getting Order
//
//	@Summary Getting Order by Code
//	@Description Getting Order by Code in detail
//	@Tags Orders
//	@Id your-mom
//	@Success 200 {string} string
//	@Router /orders/code/{orderCode} [get]
func GetHandle(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
