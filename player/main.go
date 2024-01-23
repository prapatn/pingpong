package main

import (
	"player/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	handler := handlers.NewHandler()
	app.Get("new-match", handler.NewMatch)

	app.Listen(":8888")
}
