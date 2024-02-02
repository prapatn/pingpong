package domain

import "player/pkg/models"

type MatchLogRepository interface {
	DbMigrator() (err error)
	InsertMatch() (log models.MatchLog, err error)
	GetMatchById(id string) (models.MatchLog, error)
}

type MatchLogUsecase interface {
	DbMigrator() (err error)
	InsertLog() (models.MatchLog, error)
	GetLastMatch() (matchLog models.MatchLog, err error)
	GetMatchById(id string) (matchLog models.MatchLog, err error)
}
