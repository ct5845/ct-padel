package playshared

import (
	"log/slog"
	"net/http"
	"strconv"
)

func GetPlayID(w http.ResponseWriter, r *http.Request) int {
	playID := r.PathValue("playID")
	id, err := strconv.Atoi(playID)
	if err != nil {
		slog.Error("Invalid play ID", "error", err)
		return 0
	}
	return id
}