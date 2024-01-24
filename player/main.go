package main

import (
	"player/handlers"
	"player/services"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	redisClient := initRedis()

	app := fiber.New()

	app.Use(logger.New(
		logger.Config{
			TimeZone: "Asia/Bangkok",
		},
	))

	service := services.NewService(redisClient)
	handler := handlers.NewHandler(service)
	app.Get("new-match", handler.NewMatch)
	app.Get("last-match", handler.GetLastMatch)

	app.Listen(":8888")
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
