package set

import (
	"ct-padel-s/src/features/padel/game/gamerepo"
	"ct-padel-s/src/features/padel/game/gameviews"
	"ct-padel-s/src/features/padel/match/matchrepo"
	"ct-padel-s/src/features/padel/match/matchshared"
	"ct-padel-s/src/features/padel/set/setmodel"
	"ct-padel-s/src/features/padel/set/setrepo"
	"ct-padel-s/src/features/padel/set/setshared"
	"ct-padel-s/src/features/padel/set/setviews"
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
	db := database.GetDB()

	matchID := matchshared.GetMatchID(w, r)
	if matchID == 0 {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	sets, err := setrepo.GetSetsByMatch(db, matchID)
	if err != nil {
		slog.Error("Failed to get sets", "error", err)
		http.Error(w, "Failed to get sets", http.StatusInternalServerError)
		return
	}

	set := setmodel.Set{
		MatchID:   matchID,
		SetNumber: len(sets) + 1,
	}

	if err := setrepo.CreateSet(db, &set); err != nil {
		slog.Error("Failed to create set", "error", err)
		http.Error(w, "Failed to create set", http.StatusInternalServerError)
		return
	}

	// Set HTMX redirect header and return created status
	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d", matchID, set.ID))
	w.WriteHeader(http.StatusCreated)
}

func Get(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	matchID := matchshared.GetMatchID(w, r)
	if matchID == 0 {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	setID := setshared.GetSetID(w, r)
	if setID == 0 {
		http.Error(w, "Invalid set ID", http.StatusBadRequest)
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

	breadcrumb, err := setviews.RenderBreadcrumb(match, set)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Load shared components
	title := "Match: " + match.Name() + " - Set: " + strconv.Itoa(set.SetNumber)

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

	// Get games for this set
	games, err := gamerepo.GetGamesBySet(db, set.ID)
	if err != nil {
		slog.Error("Failed to get games", "error", err, "setID", set.ID)
		http.Error(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	// Render games list
	gamesListHTML, err := gameviews.RenderGameList(games, set, match)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Load feature content and render with data
	contentHTML, err := setviews.RenderGet(set, match, games, gamesListHTML)
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

	// Check if set exists
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

	// Delete the set (this will also reorder remaining set numbers)
	if err := setrepo.DeleteSet(db, setID); err != nil {
		slog.Error("Failed to delete set", "error", err, "setID", setID)
		http.Error(w, "Failed to delete set", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d", matchID))
	w.WriteHeader(http.StatusOK)
}

func GetByMatch(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	matchID := matchshared.GetMatchID(w, r)
	if matchID == 0 {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	getByMatch(w, r, matchID, db)
}

func getByMatch(w http.ResponseWriter, r *http.Request, matchID int, db *database.DB) {
	slog.Debug("Handling", "method", r.Method, "path", r.URL.Path)

	sets, err := setrepo.GetSetsByMatch(db, matchID)
	if err != nil {
		slog.Error("Failed to get sets", "error", err, "matchID", matchID)
		http.Error(w, "Failed to get sets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sets)
	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)
}
