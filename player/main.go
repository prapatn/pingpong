package main

import (
	"fmt"
	"log"
	"os"
	"player/pkg/config"
	matchlogs "player/pkg/match_logs"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	config.GetEnv()
	redisClient := initRedis()
	// db := initMongoDB()
	db := initMySQL()

	app := fiber.New()

	app.Use(logger.New(
		logger.Config{
			TimeZone: "Asia/Bangkok",
		},
	))

	// App Repository
	matchLogReposiroty := matchlogs.NewMatchLogRepository(db)

	// App Services
	matchLogUsecase := matchlogs.NewMatchLogUsecase(matchLogReposiroty, redisClient)

	matchlogs.NewMatchLogHandler(app.Group("/"), matchLogUsecase)

	app.Listen(":8888")
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Env.Redis,
	})
}

// func initMongoDB() *mongo.Client {
// 	// monitor := &event.CommandMonitor{
// 	// 	Started: func(_ context.Context, e *event.CommandStartedEvent) {
// 	// 		log.Println("Started : ", e.Command.String())
// 	// 	},
// 	// 	Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
// 	// 		log.Println("Succeeded : ", e.Reply.String())
// 	// 	},
// 	// 	Failed: func(_ context.Context, e *event.CommandFailedEvent) {
// 	// 		log.Println("Failed : ", e.Failure)
// 	// 	},
// 	// }

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	credential := options.Credential{
// 		AuthSource: "admin",
// 		Username:   "root",
// 		Password:   "password",
// 	}
// 	connectionURI := "mongodb://127.0.0.1:27017"
// 	clientOpts := options.Client().ApplyURI(connectionURI).
// 		// SetMonitor(monitor).
// 		SetAuth(credential)
// 	client, err := mongo.Connect(ctx, clientOpts)
// 	if err != nil {
// 		log.Println("connection failed :", err)
// 	} else {
// 		log.Println("database connected")
// 	}

// 	return client
// }

func initMySQL() *gorm.DB {

	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Env.Host, config.Env.Port, config.Env.User, config.Env.Password, config.Env.DBName)

	// New logger for detailed SQL logging
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second,     // Slow SQL threshold
			LogLevel:      gormLogger.Info, // Log level
			Colorful:      true,            // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	if err != nil {
		log.Println("connection failed :", err)
	} else {
		log.Println("database connected")
	}
	return db
}
