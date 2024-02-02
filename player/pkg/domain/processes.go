package domain

import "player/pkg/models"

type ProcessRepository interface {
	DbMigrator() (err error)
	InsertProcess(process *models.Processes) error
}
