package setmodel

import "time"

type Set struct {
	ID        int       `json:"id" db:"id"`
	MatchID   int       `json:"match_id" db:"match_id"`
	SetNumber int       `json:"set_number" db:"set_number"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}