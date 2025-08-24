package playrepo

import (
	"ct-padel-s/src/features/padel/play/playmodel"
	"ct-padel-s/src/infrastructure/database"
)

func CreatePlay(db *database.DB, play *playmodel.Play) error {
	query := `INSERT INTO plays (point_id, play_number, player_id, ball_position_x, ball_position_y, result_type, hand_side, contact_type, shot_effect) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			  RETURNING id, created_at, updated_at`
	err := db.QueryRow(query, 
		play.PointID, 
		play.PlayNumber, 
		play.PlayerID, 
		play.BallPositionX, 
		play.BallPositionY, 
		play.ResultType, 
		play.HandSide, 
		play.ContactType, 
		play.ShotEffect,
	).Scan(&play.ID, &play.CreatedAt, &play.UpdatedAt)
	return err
}

func GetPlaysByPoint(db *database.DB, pointID int) ([]*playmodel.Play, error) {
	query := `SELECT id, point_id, play_number, player_id, ball_position_x, ball_position_y, result_type, hand_side, contact_type, shot_effect, created_at, updated_at 
			  FROM plays WHERE point_id = $1 ORDER BY play_number`
	rows, err := db.Query(query, pointID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plays []*playmodel.Play
	for rows.Next() {
		var play playmodel.Play
		err := rows.Scan(
			&play.ID, 
			&play.PointID, 
			&play.PlayNumber, 
			&play.PlayerID, 
			&play.BallPositionX, 
			&play.BallPositionY, 
			&play.ResultType, 
			&play.HandSide, 
			&play.ContactType, 
			&play.ShotEffect, 
			&play.CreatedAt, 
			&play.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plays = append(plays, &play)
	}
	return plays, rows.Err()
}

func GetPlay(db *database.DB, playID int) (*playmodel.Play, error) {
	query := `SELECT id, point_id, play_number, player_id, ball_position_x, ball_position_y, result_type, hand_side, contact_type, shot_effect, created_at, updated_at 
			  FROM plays WHERE id = $1`
	var play playmodel.Play
	err := db.QueryRow(query, playID).Scan(
		&play.ID, 
		&play.PointID, 
		&play.PlayNumber, 
		&play.PlayerID, 
		&play.BallPositionX, 
		&play.BallPositionY, 
		&play.ResultType, 
		&play.HandSide, 
		&play.ContactType, 
		&play.ShotEffect, 
		&play.CreatedAt, 
		&play.UpdatedAt,
	)
	return &play, err
}

func DeletePlay(db *database.DB, playID int) error {
	// First get the play to know its point_id and play_number
	play, err := GetPlay(db, playID)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the play
	_, err = tx.Exec(`DELETE FROM plays WHERE id = $1`, playID)
	if err != nil {
		return err
	}

	// Update play numbers for plays with higher numbers in the same point
	_, err = tx.Exec(`UPDATE plays SET play_number = play_number - 1 WHERE point_id = $1 AND play_number > $2`,
		play.PointID, play.PlayNumber)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func UpdatePlay(db *database.DB, play *playmodel.Play) error {
	query := `UPDATE plays SET 
				player_id = $1, 
				ball_position_x = $2, 
				ball_position_y = $3, 
				result_type = $4, 
				hand_side = $5, 
				contact_type = $6, 
				shot_effect = $7, 
				updated_at = CURRENT_TIMESTAMP
			  WHERE id = $8`
	_, err := db.Exec(query, 
		play.PlayerID, 
		play.BallPositionX, 
		play.BallPositionY, 
		play.ResultType, 
		play.HandSide, 
		play.ContactType, 
		play.ShotEffect, 
		play.ID,
	)
	return err
}

func DeleteSubsequentPlays(db *database.DB, pointID int, playNumber int) error {
	query := `DELETE FROM plays WHERE point_id = $1 AND play_number > $2`
	_, err := db.Exec(query, pointID, playNumber)
	return err
}