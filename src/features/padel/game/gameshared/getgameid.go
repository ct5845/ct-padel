package gameshared

import (
	"log/slog"
	"net/http"
	"strconv"
)

func GetGameID(w http.ResponseWriter, r *http.Request) int {
	gameID := r.PathValue("gameID")
	id, err := strconv.Atoi(gameID)
	if err != nil {
		slog.Error("Invalid game ID", "error", err)
		return 0
	}
	return id
}