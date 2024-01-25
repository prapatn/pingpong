package services

import (
	"context"
	"encoding/json"
	"player/models"
)

func (s service) GetLastMatch() (matchLog models.MatchLog, err error) {
	//GET LastMatch
	dataChk, err := s.redisClient.Get(context.Background(), "LastMatch").Result()
	if err == nil {
		err = json.Unmarshal([]byte(dataChk), &matchLog)
	}
	return matchLog, err
}
