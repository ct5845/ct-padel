package playmodel

import (
	"database/sql"
	"time"
)

type Play struct {
	ID            int            `json:"id" db:"id"`
	PointID       int            `json:"point_id" db:"point_id"`
	PlayNumber    int            `json:"play_number" db:"play_number"`
	PlayerID      sql.NullInt64  `json:"player_id" db:"player_id"`
	BallPositionX int            `json:"ball_position_x" db:"ball_position_x"`
	BallPositionY int            `json:"ball_position_y" db:"ball_position_y"`
	ResultType    sql.NullString `json:"result_type" db:"result_type"`
	HandSide      sql.NullString `json:"hand_side" db:"hand_side"`
	ContactType   sql.NullString `json:"contact_type" db:"contact_type"`
	ShotEffect    sql.NullString `json:"shot_effect" db:"shot_effect"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}