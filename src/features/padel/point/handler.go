package point

import (
	"ct-padel-s/src/features/padel/game/gamerepo"
	"ct-padel-s/src/features/padel/game/gameshared"
	"ct-padel-s/src/features/padel/match/matchrepo"
	"ct-padel-s/src/features/padel/match/matchshared"
	"ct-padel-s/src/features/padel/play/playrepo"
	"ct-padel-s/src/features/padel/play/playviews"
	"ct-padel-s/src/features/padel/point/pointmodel"
	"ct-padel-s/src/features/padel/point/pointrepo"
	"ct-padel-s/src/features/padel/point/pointshared"
	"ct-padel-s/src/features/padel/point/pointviews"
	"ct-padel-s/src/features/padel/set/setrepo"
	"ct-padel-s/src/features/padel/set/setshared"
	"ct-padel-s/src/infrastructure/database"
	"ct-padel-s/src/shared/components/footer"
	"ct-padel-s/src/shared/components/header"
	"ct-padel-s/src/shared/templates"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func Create(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Handling", "method", r.Method, "path", r.URL.Path)
	db := database.GetDB()

	matchID := matchshared.GetMatchID(w, r)
	if matchID == 0 {
		slog.Error("Invalid match ID", "matchID", matchID)
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	setID := setshared.GetSetID(w, r)
	if setID == 0 {
		slog.Error("Invalid set ID", "setID", setID)
		http.Error(w, "Invalid set ID", http.StatusBadRequest)
		return
	}

	gameID := gameshared.GetGameID(w, r)
	if gameID == 0 {
		slog.Error("Invalid game ID", "gameID", gameID)
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	points, err := pointrepo.GetPointsByGame(db, gameID)
	if err != nil {
		slog.Error("Failed to get points", "error", err)
		http.Error(w, "Failed to get points", http.StatusInternalServerError)
		return
	}

	point := pointmodel.Point{
		GameID:      gameID,
		PointNumber: len(points) + 1,
	}

	if err := pointrepo.CreatePoint(db, &point); err != nil {
		slog.Error("Failed to create point", "error", err)
		http.Error(w, "Failed to create point", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	// Set HTMX redirect header and return created status
	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d/games/%d/points/%d", matchID, setID, gameID, point.ID))
	w.WriteHeader(http.StatusCreated)
}

func Get(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Handling", "method", r.Method, "path", r.URL.Path)
	db := database.GetDB()

	matchID := matchshared.GetMatchID(w, r)
	if matchID == 0 {
		slog.Error("Invalid match ID", "matchID", matchID)
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	setID := setshared.GetSetID(w, r)
	if setID == 0 {
		slog.Error("Invalid set ID", "setID", setID)
		http.Error(w, "Invalid set ID", http.StatusBadRequest)
		return
	}

	gameID := gameshared.GetGameID(w, r)
	if gameID == 0 {
		slog.Error("Invalid game ID", "gameID", gameID)
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	pointID := pointshared.GetPointID(w, r)
	if pointID == 0 {
		slog.Error("Invalid point ID", "pointID", pointID)
		http.Error(w, "Invalid point ID", http.StatusBadRequest)
		return
	}

	match, err := matchrepo.GetMatchWithPlayers(db, matchID)
	if err != nil {
		slog.Error("Failed to get match", "error", err, "matchID", matchID)
		http.Error(w, "Failed to get match", http.StatusInternalServerError)
		return
	}

	if match == nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	set, err := setrepo.GetSet(db, setID)
	if err != nil {
		slog.Error("Failed to get set", "error", err, "setID", setID)
		http.Error(w, "Failed to get set", http.StatusInternalServerError)
		return
	}

	if set == nil {
		http.Error(w, "Set not found", http.StatusNotFound)
		return
	}

	game, err := gamerepo.GetGame(db, gameID)
	if err != nil {
		slog.Error("Failed to get game", "error", err, "gameID", gameID)
		http.Error(w, "Failed to get game", http.StatusInternalServerError)
		return
	}

	if game == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	point, err := pointrepo.GetPoint(db, pointID)
	if err != nil {
		slog.Error("Failed to get point", "error", err, "pointID", pointID)
		http.Error(w, "Failed to get point", http.StatusInternalServerError)
		return
	}

	if point == nil {
		http.Error(w, "Point not found", http.StatusNotFound)
		return
	}

	breadcrumb, err := pointviews.RenderBreadcrumb(match, set, game, point)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Load shared components
	title := "Match: " + match.Name() + " - Set: " + strconv.Itoa(set.SetNumber) + " - Game: " + strconv.Itoa(game.GameNumber) + " - Point: " + strconv.Itoa(point.PointNumber)

	headerHTML, err := header.Render(header.Data{Title: title + " - Padel Tracker", Breadcrumb: breadcrumb})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	footerHTML, err := footer.Render(footer.Data{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get plays for this point
	plays, err := playrepo.GetPlaysByPoint(db, point.ID)
	if err != nil {
		slog.Error("Failed to get plays", "error", err, "pointID", point.ID)
		http.Error(w, "Failed to get plays", http.StatusInternalServerError)
		return
	}

	// Render plays list
	playsListHTML, err := playviews.RenderPlayList(plays, point, game, set, match)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Load feature content and render with data
	contentHTML, err := pointviews.RenderGet(point, game, set, match, playsListHTML)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Compose final page
	page, err := templates.Render(templates.Data{
		Title:       title + " - Padel Tracker",
		HeaderHTML:  headerHTML,
		ContentHTML: contentHTML,
		FooterHTML:  footerHTML,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)
	io.WriteString(w, string(page))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Handling", "method", r.Method, "path", r.URL.Path)
	db := database.GetDB()

	matchID := matchshared.GetMatchID(w, r)
	if matchID == 0 {
		slog.Error("Invalid match ID", "matchID", matchID)
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	setID := setshared.GetSetID(w, r)
	if setID == 0 {
		slog.Error("Invalid set ID", "setID", setID)
		http.Error(w, "Invalid set ID", http.StatusBadRequest)
		return
	}

	gameID := gameshared.GetGameID(w, r)
	if gameID == 0 {
		slog.Error("Invalid game ID", "gameID", gameID)
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	pointID := pointshared.GetPointID(w, r)
	if pointID == 0 {
		slog.Error("Invalid point ID", "pointID", pointID)
		http.Error(w, "Invalid point ID", http.StatusBadRequest)
		return
	}

	// Check if point exists
	point, err := pointrepo.GetPoint(db, pointID)
	if err != nil {
		slog.Error("Failed to get point", "error", err, "pointID", pointID)
		http.Error(w, "Failed to get point", http.StatusInternalServerError)
		return
	}

	if point == nil {
		http.Error(w, "Point not found", http.StatusNotFound)
		return
	}

	// Delete the point (this will also reorder remaining point numbers)
	if err := pointrepo.DeletePoint(db, pointID); err != nil {
		slog.Error("Failed to delete point", "error", err, "pointID", pointID)
		http.Error(w, "Failed to delete point", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d/games/%d", matchID, setID, gameID))
	w.WriteHeader(http.StatusOK)
}

func GetByGame(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Handling", "method", r.Method, "path", r.URL.Path)
	db := database.GetDB()

	matchID := matchshared.GetMatchID(w, r)
	if matchID == 0 {
		slog.Error("Invalid match ID", "matchID", matchID)
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	setID := setshared.GetSetID(w, r)
	if setID == 0 {
		slog.Error("Invalid set ID", "setID", setID)
		http.Error(w, "Invalid set ID", http.StatusBadRequest)
		return
	}

	gameID := gameshared.GetGameID(w, r)
	if gameID == 0 {
		slog.Error("Invalid game ID", "gameID", gameID)
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	points, err := pointrepo.GetPointsByGame(db, gameID)
	if err != nil {
		slog.Error("Failed to get points", "error", err, "gameID", gameID)
		http.Error(w, "Failed to get points", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)
}