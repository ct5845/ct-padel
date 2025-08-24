package game

import (
	"ct-padel-s/src/features/padel/game/gamerepo"
	"ct-padel-s/src/features/padel/game/gameshared"
	"ct-padel-s/src/features/padel/game/gameviews"
	"ct-padel-s/src/features/padel/game/gamemodel"
	"ct-padel-s/src/features/padel/match/matchrepo"
	"ct-padel-s/src/features/padel/match/matchshared"
	"ct-padel-s/src/features/padel/point/pointrepo"
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

	games, err := gamerepo.GetGamesBySet(db, setID)
	if err != nil {
		slog.Error("Failed to get games", "error", err)
		http.Error(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	game := gamemodel.Game{
		SetID:      setID,
		GameNumber: len(games) + 1,
	}

	if err := gamerepo.CreateGame(db, &game); err != nil {
		slog.Error("Failed to create game", "error", err)
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	// Set HTMX redirect header and return created status
	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d/games/%d", matchID, setID, game.ID))
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

	breadcrumb, err := gameviews.RenderBreadcrumb(match, set, game)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Load shared components
	title := "Match: " + match.Name() + " - Set: " + strconv.Itoa(set.SetNumber) + " - Game: " + strconv.Itoa(game.GameNumber)

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

	// Get points for this game
	points, err := pointrepo.GetPointsByGame(db, game.ID)
	if err != nil {
		slog.Error("Failed to get points", "error", err, "gameID", game.ID)
		http.Error(w, "Failed to get points", http.StatusInternalServerError)
		return
	}

	// Render points list
	pointsListHTML, err := pointviews.RenderPointList(points, game, set, match)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Load feature content and render with data
	contentHTML, err := gameviews.RenderGet(game, set, match, pointsListHTML)
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

	// Check if game exists
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

	// Delete the game (this will also reorder remaining game numbers)
	if err := gamerepo.DeleteGame(db, gameID); err != nil {
		slog.Error("Failed to delete game", "error", err, "gameID", gameID)
		http.Error(w, "Failed to delete game", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d", matchID, setID))
	w.WriteHeader(http.StatusOK)
}

func GetBySet(w http.ResponseWriter, r *http.Request) {
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

	games, err := gamerepo.GetGamesBySet(db, setID)
	if err != nil {
		slog.Error("Failed to get games", "error", err, "setID", setID)
		http.Error(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)
}