package main

import (
	"log"
	matchlogs "player/pkg/match_logs"
	"player/pkg/processes"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
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
	processesReposiroty := processes.NewProcessesReposiroty(db)

	// App Services
	matchLogUsecase := matchlogs.NewMatchLogUsecase(matchLogReposiroty, processesReposiroty, redisClient)

	matchlogs.NewMatchLogHandler(app.Group("/"), matchLogUsecase)

	app.Listen(":8888")
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
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
	dial := mysql.Open("root:@tcp(localhost:3306)/local?parseTime=true")
	db, err := gorm.Open(dial, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Println("connection failed :", err)
	} else {
		log.Println("database connected")
	}
	return db
}
