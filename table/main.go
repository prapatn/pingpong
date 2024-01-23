package main

import (
	"table/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New(
		logger.Config{
			TimeZone: "Asia/Bangkok",
		},
	))
	handler := handlers.NewHandler()
	app.Get("/ping", handler.Ping)

	app.Listen(":8889")
}
