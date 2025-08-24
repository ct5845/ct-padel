package gamerepo

import (
	"ct-padel-s/src/features/padel/game/gamemodel"
	"ct-padel-s/src/infrastructure/database"
)

func CreateGame(db *database.DB, game *gamemodel.Game) error {
	query := `INSERT INTO games (set_id, game_number) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := db.QueryRow(query, game.SetID, game.GameNumber).Scan(&game.ID, &game.CreatedAt, &game.UpdatedAt)
	return err
}

func GetGamesBySet(db *database.DB, setID int) ([]*gamemodel.Game, error) {
	query := `SELECT id, set_id, game_number, created_at, updated_at FROM games WHERE set_id = $1 ORDER BY game_number`
	rows, err := db.Query(query, setID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*gamemodel.Game
	for rows.Next() {
		var game gamemodel.Game
		err := rows.Scan(&game.ID, &game.SetID, &game.GameNumber, &game.CreatedAt, &game.UpdatedAt)
		if err != nil {
			return nil, err
		}
		games = append(games, &game)
	}
	return games, rows.Err()
}

func GetGame(db *database.DB, gameID int) (*gamemodel.Game, error) {
	query := `SELECT id, set_id, game_number, created_at, updated_at FROM games WHERE id = $1`
	var game gamemodel.Game
	err := db.QueryRow(query, gameID).Scan(&game.ID, &game.SetID, &game.GameNumber, &game.CreatedAt, &game.UpdatedAt)
	return &game, err
}

func DeleteGame(db *database.DB, gameID int) error {
	// First get the game to know its set_id and game_number
	game, err := GetGame(db, gameID)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the game
	_, err = tx.Exec(`DELETE FROM games WHERE id = $1`, gameID)
	if err != nil {
		return err
	}

	// Update game numbers for games with higher numbers in the same set
	_, err = tx.Exec(`UPDATE games SET game_number = game_number - 1 WHERE set_id = $1 AND game_number > $2`,
		game.SetID, game.GameNumber)
	if err != nil {
		return err
	}

	return tx.Commit()
}
