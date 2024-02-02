package repositories

import (
	"player/models"
)

func (r repository) GetMatchById(id string) (matchLog models.MatchLog, err error) {
	// Mongo
	// filter := bson.D{{Key: "_id", Value: id}}

	// err = r.db.Database("local").Collection("match_log").FindOne(r.ctx, filter).Decode(&matchLog)
	// if err == mongo.ErrNoDocuments {
	// 	return matchLog, errors.New("Match not found")
	// }

	//MySQL

	err = r.db.Table("match_logs").Preload("Processes").Where(id).First(&matchLog).Error

	return matchLog, err

}
