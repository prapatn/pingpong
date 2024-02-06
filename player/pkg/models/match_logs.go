package models

import "time"

type MatchLog struct {
	ID          uint      `json:"id"`
	MatchNumber string    `json:"match_number"`
	Player      string    `json:"player"`
	Turn        int       `json:"turn"`
	BallPower   int       `json:"ball_power"`
	Time        time.Time `json:"time"`
}
