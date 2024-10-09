package main

import (
	"fmt"
	"log"
	"net/http"
	"order_service/composer"
	"order_service/config"
	"os"
	"os/signal"
	"syscall"

	_ "order_service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	cfg := config.NewConfig()
	pg := config.ConnectToPostgres(cfg)
	rd := config.ConnectToRedis(cfg)
	s3Client := config.ConnectToAWS(cfg)

	defer pg.Close()
	defer rd.Close()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use("/static", filesystem.New(
		filesystem.Config{
			Root:   http.Dir(fmt.Sprintf("%s/storage", pwd)),
			Browse: true,
			MaxAge: 3600,
		}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	composer.SetUpRoutes(app.Group("/v1"), cfg, pg, rd, s3Client)

	go func() {
		log.Println("App runnning, Ctrl + C to shut down")
		log.Fatalln(app.Listen("0.0.0.0:8080"))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
}
