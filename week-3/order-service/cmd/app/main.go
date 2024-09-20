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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.NewConfig()
	pg := config.ConnectToPostgres(cfg)
	rd := config.ConnectToRedis(cfg)
	defer pg.Close()
	defer rd.Close()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use("/static", filesystem.New(
		filesystem.Config{
			Root:   http.Dir(fmt.Sprintf("%s/storage", pwd)),
			Browse: true,
			MaxAge: 3600,
		}))

	composer.SetUpRoutes(app.Group("/v1"), cfg, pg, rd)

	go func() {
		log.Println("App runnning, Ctrl + C to shut down")
		log.Fatalln(app.Listen("0.0.0.0:8080"))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
}
