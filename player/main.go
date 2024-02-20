package main

import (
	"fmt"
	"log"
	"os"
	matchlogs "player/pkg/match_logs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// redisClient := initRedis()
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
	matchLogUsecase := matchlogs.NewMatchLogUsecase(matchLogReposiroty, nil)

	matchlogs.NewMatchLogHandler(app.Group("/"), matchLogUsecase)

	app.Listen(":8888")
}

// func initRedis() *redis.Client {
// 	return redis.NewClient(&redis.Options{
// 		Addr: "10.11.232.171:6379",
// 	})
// }

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
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read database configuration from .env file
	host := os.Getenv("DB_HOST")
	// port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Configure your MySQL database details here
	dsn := fmt.Sprintf("%s:%s@unix(%s)/%s", user, password, host, dbname)

	// New logger for detailed SQL logging
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second,     // Slow SQL threshold
			LogLevel:      gormLogger.Info, // Log level
			Colorful:      true,            // Enable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
