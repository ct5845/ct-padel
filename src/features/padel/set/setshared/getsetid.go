package setshared

import (
	"log/slog"
	"net/http"
	"strconv"
)

func GetSetID(w http.ResponseWriter, r *http.Request) int {
	setID := r.PathValue("setID")
	id, err := strconv.Atoi(setID)
	if err != nil {
		slog.Error("Invalid set ID", "error", err)
		return 0
	}
	return id
}
