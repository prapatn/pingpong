package processes

import (
	"player/pkg/domain"
	"player/pkg/models"

	"gorm.io/gorm"
)

type processesReposiroty struct {
	DB *gorm.DB
}

func NewProcessesReposiroty(db *gorm.DB) domain.ProcessRepository {
	processesReposiroty := &processesReposiroty{
		DB: db,
	}
	processesReposiroty.DbMigrator()
	return processesReposiroty
}

func (r *processesReposiroty) DbMigrator() (err error) {
	err = r.DB.AutoMigrate(&models.Processes{})
	return
}

func (r *processesReposiroty) InsertProcess(process *models.Processes) error {
	return r.DB.Create(&process).Error
}
