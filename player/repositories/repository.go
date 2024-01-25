package repositories

import (
	"context"
	"player/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertMatch(log models.MatchLog) error
	GetMatchById(id string) (models.MatchLog, error)
}

type repository struct {
	db  *mongo.Client
	ctx context.Context
}

func NewRepository(db *mongo.Client) Repository {
	ctx := context.TODO()
	return repository{db: db, ctx: ctx}
}
