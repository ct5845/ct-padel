package pointmodel

import "time"

type Point struct {
	ID          int       `json:"id" db:"id"`
	GameID      int       `json:"game_id" db:"game_id"`
	PointNumber int       `json:"point_number" db:"point_number"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}