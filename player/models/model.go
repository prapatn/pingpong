package models

type MatchLog struct {
	ID      string    `bson:"_id"`
	Process []Process `bson:"process"`
}

type Process struct {
	Player    string `bson:"player"`
	Turn      int    `bson:"turn"`
	BallPower int    `bson:"ball_power"`
	Time      string `bson:"time"`
}
