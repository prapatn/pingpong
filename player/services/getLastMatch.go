package services

import (
	"context"
	"encoding/json"
)

func (s service) GetLastMatch() (matchLog MatchLog, err error) {
	//GET LastMatch
	dataChk, err := s.redisClient.Get(context.Background(), "LastMatch").Result()
	if err == nil {
		err = json.Unmarshal([]byte(dataChk), &matchLog)
	}
	return matchLog, err
}
