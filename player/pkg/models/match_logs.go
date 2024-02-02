package models

type MatchLog struct {
	ID        uint        `json:"id"`
	Processes []Processes `json:"processes" gorm:"foreignKey:MatchLogID"`
}
