package matchmodel

import (
	"ct-padel-s/src/features/padel/player/playermodel"
	"fmt"
	"time"
)

type Match struct {
	ID             int       `json:"id" db:"id"`
	Team1Player1ID int       `json:"team1_player1_id" db:"team1_player1_id"`
	Team1Player2ID int       `json:"team1_player2_id" db:"team1_player2_id"`
	Team2Player1ID int       `json:"team2_player1_id" db:"team2_player1_id"`
	Team2Player2ID int       `json:"team2_player2_id" db:"team2_player2_id"`
	MatchDate      time.Time `json:"match_date" db:"match_date"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type MatchWithPlayers struct {
	Match
	Team1Player1 playermodel.Player `json:"team1_player1"`
	Team1Player2 playermodel.Player `json:"team1_player2"`
	Team2Player1 playermodel.Player `json:"team2_player1"`
	Team2Player2 playermodel.Player `json:"team2_player2"`
}

func (m *MatchWithPlayers) Name() string {
	return fmt.Sprintf("%s | %s & %s vs %s & %s",
		m.MatchDate.Format("Mon, 02 Jan 15:04:05"),
		m.Team1Player1.Name, m.Team1Player2.Name, m.Team2Player1.Name, m.Team2Player2.Name)
}
