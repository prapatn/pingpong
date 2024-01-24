package services

import (
	"github.com/go-redis/redis/v8"
)

type Services interface {
	InsertLog() error
	GetLastMatch() (matchLog MatchLog, err error)
}

type service struct {
	redisClient *redis.Client
}

func NewService(redisClient *redis.Client) Services {
	return service{redisClient: redisClient}
}

type MatchLog struct {
	Number  int       `bson:"number"`
	Process []Process `bson:"process"`
}

type Process struct {
	Player    string `bson:"player"`
	Turn      int    `bson:"turn"`
	BallPower int    `bson:"ball_power"`
	Time      string `bson:"time"`
}
