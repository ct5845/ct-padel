package gamemodel

import "time"

type Game struct {
	ID         int       `json:"id" db:"id"`
	SetID      int       `json:"set_id" db:"set_id"`
	GameNumber int       `json:"game_number" db:"game_number"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
