package matchlogs

import (
	"player/pkg/domain"
	"player/pkg/models"

	"gorm.io/gorm"
)

type matchLogReposiroty struct {
	DB *gorm.DB
}

func NewMatchLogRepository(db *gorm.DB) domain.MatchLogRepository {
	matchLogReposiroty := &matchLogReposiroty{
		DB: db,
	}
	matchLogReposiroty.DbMigrator()
	return matchLogReposiroty
}

func (r *matchLogReposiroty) DbMigrator() (err error) {
	err = r.DB.AutoMigrate(&models.MatchLog{})
	return
}

func (r matchLogReposiroty) GetMatchByMacthNumber(number string) (matchLog []models.MatchLog, err error) {
	// Mongo
	// filter := bson.D{{Key: "_id", Value: id}}

	// err = r.db.Database("local").Collection("match_log").FindOne(r.ctx, filter).Decode(&matchLog)
	// if err == mongo.ErrNoDocuments {
	// 	return matchLog, errors.New("Match not found")
	// }

	//MySQL

	err = r.DB.Table("match_logs").Where("match_number = ?", number).Order("id").Find(&matchLog).Error

	return matchLog, err

}

func (r matchLogReposiroty) InsertMatch(log models.MatchLog) (id int, err error) {
	//Mongo

	// result, err := r.db.Database("local").Collection("match_log").InsertOne(r.ctx, log)
	// if err != nil {
	// 	return errors.New(fmt.Sprintf("error create project : %s", err))
	// }
	// if result.InsertedID == "" {
	// 	return errors.New(fmt.Sprintf("no document was add : %s", err))
	// }

	//MySQL

	err = r.DB.Create(&log).Error

	return int(log.ID), err
}
