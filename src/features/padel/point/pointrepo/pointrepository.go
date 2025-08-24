package pointrepo

import (
	"ct-padel-s/src/features/padel/point/pointmodel"
	"ct-padel-s/src/infrastructure/database"
)

func CreatePoint(db *database.DB, point *pointmodel.Point) error {
	query := `INSERT INTO points (game_id, point_number) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := db.QueryRow(query, point.GameID, point.PointNumber).Scan(&point.ID, &point.CreatedAt, &point.UpdatedAt)
	return err
}

func GetPointsByGame(db *database.DB, gameID int) ([]*pointmodel.Point, error) {
	query := `SELECT id, game_id, point_number, created_at, updated_at FROM points WHERE game_id = $1 ORDER BY point_number`
	rows, err := db.Query(query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*pointmodel.Point
	for rows.Next() {
		var point pointmodel.Point
		err := rows.Scan(&point.ID, &point.GameID, &point.PointNumber, &point.CreatedAt, &point.UpdatedAt)
		if err != nil {
			return nil, err
		}
		points = append(points, &point)
	}
	return points, rows.Err()
}

func GetPoint(db *database.DB, pointID int) (*pointmodel.Point, error) {
	query := `SELECT id, game_id, point_number, created_at, updated_at FROM points WHERE id = $1`
	var point pointmodel.Point
	err := db.QueryRow(query, pointID).Scan(&point.ID, &point.GameID, &point.PointNumber, &point.CreatedAt, &point.UpdatedAt)
	return &point, err
}

func DeletePoint(db *database.DB, pointID int) error {
	// First get the point to know its game_id and point_number
	point, err := GetPoint(db, pointID)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the point
	_, err = tx.Exec(`DELETE FROM points WHERE id = $1`, pointID)
	if err != nil {
		return err
	}

	// Update point numbers for points with higher numbers in the same game
	_, err = tx.Exec(`UPDATE points SET point_number = point_number - 1 WHERE game_id = $1 AND point_number > $2`,
		point.GameID, point.PointNumber)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CreateNextPoint(db *database.DB, gameID int) (*pointmodel.Point, error) {
	// Get existing points for this game
	points, err := GetPointsByGame(db, gameID)
	if err != nil {
		return nil, err
	}

	// Create new point with next number
	point := &pointmodel.Point{
		GameID:      gameID,
		PointNumber: len(points) + 1,
	}

	if err := CreatePoint(db, point); err != nil {
		return nil, err
	}

	return point, nil
}