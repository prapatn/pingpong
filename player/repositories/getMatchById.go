package repositories

import (
	"errors"
	"player/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r repository) GetMatchById(id string) (matchLog models.MatchLog, err error) {
	filter := bson.D{{Key: "_id", Value: id}}

	err = r.db.Database("local").Collection("match_log").FindOne(r.ctx, filter).Decode(&matchLog)
	if err == mongo.ErrNoDocuments {
		return matchLog, errors.New("Match not found")
	}

	return matchLog, err

}
