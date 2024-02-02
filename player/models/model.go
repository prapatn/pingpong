package models

import (
	"time"
)

type MatchLog struct {
	ID        uint        `json:"id"`
	Processes []Processes `json:"processes" gorm:"foreignKey:MatchLogID"`
}

type Processes struct {
	ID         uint      `json:"id"`
	MatchLogID uint      `json:"match_log_id" gorm:"index"`
	Player     string    `json:"player"`
	Turn       int       `json:"turn"`
	BallPower  int       `json:"ball_power"`
	Time       time.Time `json:"time"`
}
