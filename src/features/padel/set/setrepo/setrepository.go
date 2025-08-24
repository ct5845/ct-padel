package setrepo

import (
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/infrastructure/database"
)

func CreateSet(db *database.DB, set *setmodel.Set) error {
	query := `INSERT INTO sets (match_id, set_number) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := db.QueryRow(query, set.MatchID, set.SetNumber).Scan(&set.ID, &set.CreatedAt, &set.UpdatedAt)
	return err
}

func GetAllSets(db *database.DB) ([]*setmodel.Set, error) {
	query := `SELECT id, match_id, set_number, created_at, updated_at FROM sets ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []*setmodel.Set
	for rows.Next() {
		var set setmodel.Set
		err := rows.Scan(&set.ID, &set.MatchID, &set.SetNumber, &set.CreatedAt, &set.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sets = append(sets, &set)
	}
	return sets, rows.Err()
}

func GetSetsByMatch(db *database.DB, matchID int) ([]*setmodel.Set, error) {
	query := `SELECT id, match_id, set_number, created_at, updated_at FROM sets WHERE match_id = $1 ORDER BY set_number`
	rows, err := db.Query(query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []*setmodel.Set
	for rows.Next() {
		var set setmodel.Set
		err := rows.Scan(&set.ID, &set.MatchID, &set.SetNumber, &set.CreatedAt, &set.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sets = append(sets, &set)
	}
	return sets, rows.Err()
}

func GetSet(db *database.DB, setID int) (*setmodel.Set, error) {
	query := `SELECT id, match_id, set_number, created_at, updated_at FROM sets WHERE id = $1`
	var set setmodel.Set
	err := db.QueryRow(query, setID).Scan(&set.ID, &set.MatchID, &set.SetNumber, &set.CreatedAt, &set.UpdatedAt)
	return &set, err
}

func DeleteSet(db *database.DB, setID int) error {
	// First get the set to know its match_id and set_number
	set, err := GetSet(db, setID)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the set
	_, err = tx.Exec(`DELETE FROM sets WHERE id = $1`, setID)
	if err != nil {
		return err
	}

	// Update set numbers for sets with higher numbers in the same match
	_, err = tx.Exec(`UPDATE sets SET set_number = set_number - 1 WHERE match_id = $1 AND set_number > $2`, 
		set.MatchID, set.SetNumber)
	if err != nil {
		return err
	}

	return tx.Commit()
}
