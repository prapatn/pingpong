package services

import (
	"player/models"
)

func (s service) GetMatchById(id string) (matchLog models.MatchLog, err error) {
	return s.repo.GetMatchById(id)
}
