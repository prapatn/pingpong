package repositories

import (
	"player/models"
)

func (r repository) InsertMatch() (log models.MatchLog, err error) {
	//Mongo

	// result, err := r.db.Database("local").Collection("match_log").InsertOne(r.ctx, log)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("error create project : %s", err))
	// }
	// if result.InsertedID == "" {
	// 	return errors.New(fmt.Sprintf("no document was add : %s", err))
	// }

	//MySQL

	err = r.db.Create(&log).Error

	return log, err
}

func (r repository) InsertProcess(process *models.Processes) error {
	return r.db.Create(&process).Error
}
