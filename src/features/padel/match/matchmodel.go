package match

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
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func (m *Match) Name() string {
	return fmt.Sprintf("%d & %d vs %d & %d", m.Team1Player1ID, m.Team1Player2ID, m.Team2Player1ID, m.Team2Player2ID)
}

type MatchWithPlayers struct {
	Match
	Team1Player1 playermodel.Player `json:"team1_player1"`
	Team1Player2 playermodel.Player `json:"team1_player2"`
	Team2Player1 playermodel.Player `json:"team2_player1"`
	Team2Player2 playermodel.Player `json:"team2_player2"`
}

func (m *MatchWithPlayers) Name() string {
	return fmt.Sprintf("%s & %s vs %s & %s", m.Team1Player1.Name, m.Team1Player2.Name, m.Team2Player1.Name, m.Team2Player2.Name)
}
