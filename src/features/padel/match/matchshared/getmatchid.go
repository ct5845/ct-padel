package matchshared

import (
	"log/slog"
	"net/http"
	"strconv"
)

func GetMatchID(w http.ResponseWriter, r *http.Request) int {
	matchID := r.PathValue("matchID")
	id, err := strconv.Atoi(matchID)
	if err != nil {
		slog.Error("Invalid match ID", "error", err)
		return 0
	}
	return id
}
