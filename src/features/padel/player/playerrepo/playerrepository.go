package playerrepo

import (
	"ct-padel-s/src/features/padel/player/playermodel"
	"ct-padel-s/src/infrastructure/database"
	"database/sql"
)

func CreatePlayer(db *database.DB, player *playermodel.Player) error {
	query := `INSERT INTO players (name) VALUES ($1) RETURNING id, created_at`
	err := db.QueryRow(query, player.Name).Scan(&player.ID, &player.CreatedAt)
	return err
}

func GetPlayer(db *database.DB, id int) (*playermodel.Player, error) {
	player := &playermodel.Player{}
	query := `SELECT id, name, created_at FROM players WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&player.ID, &player.Name, &player.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return player, err
}

func GetAllPlayers(db *database.DB) ([]*playermodel.Player, error) {
	query := `SELECT id, name, created_at FROM players ORDER BY name`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*playermodel.Player
	for rows.Next() {
		var player playermodel.Player
		err := rows.Scan(&player.ID, &player.Name, &player.CreatedAt)
		if err != nil {
			return nil, err
		}
		players = append(players, &player)
	}
	return players, rows.Err()
}

func DeletePlayer(db *database.DB, id int) error {
	query := `DELETE FROM players WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}