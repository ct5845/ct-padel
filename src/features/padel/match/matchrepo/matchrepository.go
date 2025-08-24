package matchrepo

import (
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/infrastructure/database"
	"database/sql"
)

func CreateMatch(db *database.DB, match *matchmodel.Match) error {
	query := `INSERT INTO matches (team1_player1_id, team1_player2_id, team2_player1_id, team2_player2_id)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, match_date, created_at, updated_at`
	err := db.QueryRow(query,
		match.Team1Player1ID,
		match.Team1Player2ID,
		match.Team2Player1ID,
		match.Team2Player2ID).Scan(&match.ID, &match.MatchDate, &match.CreatedAt, &match.UpdatedAt)
	return err
}

func GetMatch(db *database.DB, id int) (*matchmodel.Match, error) {
	match := &matchmodel.Match{}
	query := `SELECT id, team1_player1_id, team1_player2_id, team2_player1_id, team2_player2_id, match_date, created_at, updated_at
			  FROM matches WHERE id = $1`
	err := db.QueryRow(query, id).Scan(
		&match.ID,
		&match.Team1Player1ID,
		&match.Team1Player2ID,
		&match.Team2Player1ID,
		&match.Team2Player2ID,
		&match.MatchDate,
		&match.CreatedAt,
		&match.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return match, err
}

func GetAllMatches(db *database.DB) ([]matchmodel.MatchWithPlayers, error) {
	query := `SELECT
		m.id, m.team1_player1_id, m.team1_player2_id, m.team2_player1_id, m.team2_player2_id, m.match_date, m.created_at, m.updated_at,
		p1.id, p1.name, p1.created_at,
		p2.id, p2.name, p2.created_at,
		p3.id, p3.name, p3.created_at,
		p4.id, p4.name, p4.created_at
	FROM matches m
	JOIN players p1 ON m.team1_player1_id = p1.id
	JOIN players p2 ON m.team1_player2_id = p2.id
	JOIN players p3 ON m.team2_player1_id = p3.id
	JOIN players p4 ON m.team2_player2_id = p4.id
	ORDER BY m.match_date DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []matchmodel.MatchWithPlayers
	for rows.Next() {
		var match matchmodel.MatchWithPlayers
		err := rows.Scan(
			&match.ID, &match.Team1Player1ID, &match.Team1Player2ID,
			&match.Team2Player1ID, &match.Team2Player2ID, &match.MatchDate, &match.CreatedAt, &match.UpdatedAt,
			&match.Team1Player1.ID, &match.Team1Player1.Name, &match.Team1Player1.CreatedAt,
			&match.Team1Player2.ID, &match.Team1Player2.Name, &match.Team1Player2.CreatedAt,
			&match.Team2Player1.ID, &match.Team2Player1.Name, &match.Team2Player1.CreatedAt,
			&match.Team2Player2.ID, &match.Team2Player2.Name, &match.Team2Player2.CreatedAt)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, rows.Err()
}

func GetMatchWithPlayers(db *database.DB, id int) (*matchmodel.MatchWithPlayers, error) {
	match := &matchmodel.MatchWithPlayers{}
	query := `SELECT
		m.id, m.team1_player1_id, m.team1_player2_id, m.team2_player1_id, m.team2_player2_id, m.match_date, m.created_at, m.updated_at,
		p1.id, p1.name, p1.created_at,
		p2.id, p2.name, p2.created_at,
		p3.id, p3.name, p3.created_at,
		p4.id, p4.name, p4.created_at
	FROM matches m
	JOIN players p1 ON m.team1_player1_id = p1.id
	JOIN players p2 ON m.team1_player2_id = p2.id
	JOIN players p3 ON m.team2_player1_id = p3.id
	JOIN players p4 ON m.team2_player2_id = p4.id
	WHERE m.id = $1`

	err := db.QueryRow(query, id).Scan(
		&match.ID, &match.Team1Player1ID, &match.Team1Player2ID,
		&match.Team2Player1ID, &match.Team2Player2ID, &match.MatchDate, &match.CreatedAt, &match.UpdatedAt,
		&match.Team1Player1.ID, &match.Team1Player1.Name, &match.Team1Player1.CreatedAt,
		&match.Team1Player2.ID, &match.Team1Player2.Name, &match.Team1Player2.CreatedAt,
		&match.Team2Player1.ID, &match.Team2Player1.Name, &match.Team2Player1.CreatedAt,
		&match.Team2Player2.ID, &match.Team2Player2.Name, &match.Team2Player2.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return match, err
}

func DeleteMatch(db *database.DB, id int) error {
	query := `DELETE FROM matches WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
