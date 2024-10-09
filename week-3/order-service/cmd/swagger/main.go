package main

import (
	_ "order_service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Order Service API
// @version 1.0
// @description An order server handles order requests

// @contact.name partridge
// @contact.email anhduc130703@gmail.com

// @license.name MIT

// @host localhost:8080
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**
func main() {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Listen(":8080")
}
