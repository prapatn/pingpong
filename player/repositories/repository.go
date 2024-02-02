package repositories

import (
	"player/models"

	"gorm.io/gorm"
)

type Repository interface {
	InsertMatch() (log models.MatchLog, err error)
	GetMatchById(id string) (models.MatchLog, error)
	InsertProcess(process *models.Processes) error
}

type repository struct {
	db *gorm.DB
	// ctx context.Context
}

func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&models.MatchLog{})
	db.AutoMigrate(&models.Processes{})
	return repository{db: db}
}
