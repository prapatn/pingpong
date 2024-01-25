package main

import (
	"context"
	"log"
	"player/handlers"
	"player/repositories"
	"player/services"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	redisClient := initRedis()
	db := initDB()

	app := fiber.New()

	app.Use(logger.New(
		logger.Config{
			TimeZone: "Asia/Bangkok",
		},
	))

	repository := repositories.NewRepository(db)
	service := services.NewService(redisClient, repository)
	handler := handlers.NewHandler(service)
	app.Get("new-match", handler.NewMatch)
	app.Get("match", handler.GetLastMatch)

	app.Listen(":8888")
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func initDB() *mongo.Client {
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			log.Println("Started : ", e.Command.String())
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			log.Println("Succeeded : ", e.Reply.String())
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			log.Println("Failed : ", e.Failure)
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	credential := options.Credential{
		AuthSource: "admin",
		Username:   "root",
		Password:   "password",
	}
	connectionURI := "mongodb://127.0.0.1:27017"
	clientOpts := options.Client().ApplyURI(connectionURI).SetMonitor(monitor).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Println("connection failed :", err)
	} else {
		log.Println("database connected")
	}

	return client
}
