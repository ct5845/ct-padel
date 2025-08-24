package match

import (
	"ct-padel-s/src/features/padel/match/matchmodel"
	"ct-padel-s/src/features/padel/match/matchrepo"
	"ct-padel-s/src/features/padel/match/matchviews"
	"ct-padel-s/src/features/padel/player/playermodel"
	"ct-padel-s/src/features/padel/player/playerrepo"
	"ct-padel-s/src/features/padel/set/setrepo"
	"ct-padel-s/src/infrastructure/database"
	"ct-padel-s/src/shared/components/footer"
	"ct-padel-s/src/shared/components/header"
	"ct-padel-s/src/shared/templates"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Handling", "method", r.Method, "path", r.URL.Path)
	db := database.GetDB()

	matches, err := matchrepo.GetAllMatches(db)
	if err != nil {
		slog.Error("Failed to get matches", "error", err)
		http.Error(w, "Failed to get matches", http.StatusInternalServerError)
		return
	}

	breadcrumb, err := matchviews.RenderGetAllBreadcrumb()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Load shared components
	headerHTML, err := header.Render(header.Data{Title: "Matches - Padel Tracker", Breadcrumb: breadcrumb})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	footerHTML, err := footer.Render(footer.Data{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Load feature content and render with data
	contentHTML, err := matchviews.RenderGetAll(matches)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Compose final page
	page, err := templates.Render(templates.Data{
		Title:       "Matches - CT Padel Tracker",
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

func Create(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	p1 := playermodel.Player{Name: "P1"}
	p2 := playermodel.Player{Name: "P2"}
	p3 := playermodel.Player{Name: "P3"}
	p4 := playermodel.Player{Name: "P4"}

	if err := playerrepo.CreatePlayer(db, &p1); err != nil {
		slog.Error("Failed to create player", "error", err)
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
		return
	}
	if err := playerrepo.CreatePlayer(db, &p2); err != nil {
		slog.Error("Failed to create player", "error", err)
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
		return
	}
	if err := playerrepo.CreatePlayer(db, &p3); err != nil {
		slog.Error("Failed to create player", "error", err)
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
		return
	}
	if err := playerrepo.CreatePlayer(db, &p4); err != nil {
		slog.Error("Failed to create player", "error", err)
		http.Error(w, "Failed to create player", http.StatusInternalServerError)
		return
	}

	// Create a basic match with default players (we'll need to update this later)
	match := matchmodel.Match{
		Team1Player1ID: p1.ID,
		Team1Player2ID: p2.ID,
		Team2Player1ID: p3.ID,
		Team2Player2ID: p4.ID,
	}

	if err := matchrepo.CreateMatch(db, &match); err != nil {
		slog.Error("Failed to create match", "error", err)
		http.Error(w, "Failed to create match", http.StatusInternalServerError)
		return
	}

	// Set HTMX redirect header and return created status
	w.Header().Set("HX-Redirect", "/matches/"+strconv.Itoa(match.ID))
	w.WriteHeader(http.StatusCreated)
}

func Get(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Handling", "method", r.Method, "path", r.URL.Path)
	db := database.GetDB()

	matchID := r.PathValue("matchID")
	id, err := strconv.Atoi(matchID)
	if err != nil {
		slog.Error("Invalid match ID", "error", err)
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	match, err := matchrepo.GetMatchWithPlayers(db, id)
	if err != nil {
		slog.Error("Failed to get match", "error", err, "id", id)
		http.Error(w, "Failed to get match", http.StatusInternalServerError)
		return
	}

	if match == nil {
		slog.Error("Match not found", "id", id)
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	sets, err := setrepo.GetSetsByMatch(db, match.ID)
	if err != nil {
		slog.Error("Failed to get sets", "error", err, "matchID", match.ID)
		http.Error(w, "Failed to get sets", http.StatusInternalServerError)
		return
	}

	// Load shared components
	title := "Match: " + match.Name()

	breadcrumb, err := matchviews.RenderGetBreadcrumb(match)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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

	// Load feature content and render with data
	contentHTML, err := matchviews.RenderGet(match, sets)
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

	matchID := r.PathValue("matchID")
	id, err := strconv.Atoi(matchID)

	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	// Check if match exists
	match, err := matchrepo.GetMatch(db, id)
	if err != nil {
		slog.Error("Failed to get match", "error", err, "id", id)
		http.Error(w, "Failed to get match", http.StatusInternalServerError)
		return
	}

	if match == nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	// Delete the match
	if err := matchrepo.DeleteMatch(db, id); err != nil {
		slog.Error("Failed to delete match", "error", err, "id", id)
		http.Error(w, "Failed to delete match", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("HX-Redirect", "/matches")
	w.WriteHeader(http.StatusOK)
}
