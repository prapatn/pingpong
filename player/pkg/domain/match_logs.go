package domain

import "player/pkg/models"

type MatchLogRepository interface {
	DbMigrator() (err error)
	InsertMatch(log models.MatchLog) (id int, err error)
	GetMatchByMacthNumber(number string) ([]models.MatchLog, error)
}

type MatchLogUsecase interface {
	DbMigrator() (err error)
	InsertLog() ([]models.MatchLog, error)
	GetLastMatch() (matchLogs []models.MatchLog, err error)
	GetMatchByMacthNumber(number string) (matchLog []models.MatchLog, err error)
}
