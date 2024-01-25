package repositories

import (
	"errors"
	"fmt"
	"player/models"
)

func (r repository) InsertMatch(log models.MatchLog) error {
	result, err := r.db.Database("local").Collection("match_log").InsertOne(r.ctx, log)
	if err != nil {
		return errors.New(fmt.Sprintf("error create project : %s", err))
	}
	if result.InsertedID == "" {
		return errors.New(fmt.Sprintf("no document was add : %s", err))
	}
	return nil
}
