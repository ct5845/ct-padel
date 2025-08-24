package pointshared

import (
	"log/slog"
	"net/http"
	"strconv"
)

func GetPointID(w http.ResponseWriter, r *http.Request) int {
	pointID := r.PathValue("pointID")
	id, err := strconv.Atoi(pointID)
	if err != nil {
		slog.Error("Invalid point ID", "error", err)
		return 0
	}
	return id
}