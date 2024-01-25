package services

import (
	"player/models"
	"player/repositories"

	"github.com/go-redis/redis/v8"
)

type Services interface {
	InsertLog() (models.MatchLog, error)
	GetLastMatch() (matchLog models.MatchLog, err error)
	GetMatchById(id string) (matchLog models.MatchLog, err error)
}

type service struct {
	redisClient *redis.Client
	repo        repositories.Repository
}

func NewService(redisClient *redis.Client, repo repositories.Repository) Services {
	return service{redisClient: redisClient, repo: repo}
}
