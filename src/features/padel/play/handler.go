package play

import (
	"ct-padel-s/src/features/padel/game/gamerepo"
	"ct-padel-s/src/features/padel/game/gameshared"
	"ct-padel-s/src/features/padel/match/matchrepo"
	"ct-padel-s/src/features/padel/match/matchshared"
	"ct-padel-s/src/features/padel/play/playmodel"
	"ct-padel-s/src/features/padel/play/playrepo"
	"ct-padel-s/src/features/padel/play/playshared"
	"ct-padel-s/src/features/padel/play/playviews"
	"ct-padel-s/src/features/padel/point/pointrepo"
	"ct-padel-s/src/features/padel/point/pointshared"
	"ct-padel-s/src/features/padel/set/setrepo"
	"ct-padel-s/src/features/padel/set/setshared"
	"ct-padel-s/src/infrastructure/database"
	"ct-padel-s/src/shared/components/footer"
	"ct-padel-s/src/shared/components/header"
	"ct-padel-s/src/shared/templates"
	"database/sql"
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

	pointID := pointshared.GetPointID(w, r)
	if pointID == 0 {
		slog.Error("Invalid point ID", "pointID", pointID)
		http.Error(w, "Invalid point ID", http.StatusBadRequest)
		return
	}

	plays, err := playrepo.GetPlaysByPoint(db, pointID)
	if err != nil {
		slog.Error("Failed to get plays", "error", err)
		http.Error(w, "Failed to get plays", http.StatusInternalServerError)
		return
	}

	// Check if the point has already ended (last play has a result_type)
	if len(plays) > 0 {
		lastPlay := plays[len(plays)-1]
		if lastPlay.ResultType.Valid && lastPlay.ResultType.String != "" {
			slog.Error("Cannot create play after point has ended", "pointID", pointID, "lastPlayResult", lastPlay.ResultType.String)
			http.Error(w, "Cannot create play after point has ended", http.StatusBadRequest)
			return
		}
	}

	play := playmodel.Play{
		PointID:       pointID,
		PlayNumber:    len(plays) + 1,
		PlayerID:      sql.NullInt64{Valid: false},
		BallPositionX: 0,
		BallPositionY: 0,
		ResultType:    sql.NullString{Valid: false},
		HandSide:      sql.NullString{Valid: false},
		ContactType:   sql.NullString{Valid: false},
		ShotEffect:    sql.NullString{Valid: false},
	}

	if err := playrepo.CreatePlay(db, &play); err != nil {
		slog.Error("Failed to create play", "error", err)
		http.Error(w, "Failed to create play", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	// Set HTMX redirect header and return created status
	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d/games/%d/points/%d/plays/%d", matchID, setID, gameID, pointID, play.ID))
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

	playID := playshared.GetPlayID(w, r)
	if playID == 0 {
		slog.Error("Invalid play ID", "playID", playID)
		http.Error(w, "Invalid play ID", http.StatusBadRequest)
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

	play, err := playrepo.GetPlay(db, playID)
	if err != nil {
		slog.Error("Failed to get play", "error", err, "playID", playID)
		http.Error(w, "Failed to get play", http.StatusInternalServerError)
		return
	}

	if play == nil {
		http.Error(w, "Play not found", http.StatusNotFound)
		return
	}

	breadcrumb, err := playviews.RenderBreadcrumb(match, set, game, point, play)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Load shared components
	title := "Match: " + match.Name() + " - Set: " + strconv.Itoa(set.SetNumber) + " - Game: " + strconv.Itoa(game.GameNumber) + " - Point: " + strconv.Itoa(point.PointNumber) + " - Play: " + strconv.Itoa(play.PlayNumber)

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
	contentHTML, err := playviews.RenderGet(play, point, game, set, match)
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

	playID := playshared.GetPlayID(w, r)
	if playID == 0 {
		slog.Error("Invalid play ID", "playID", playID)
		http.Error(w, "Invalid play ID", http.StatusBadRequest)
		return
	}

	// Check if play exists
	play, err := playrepo.GetPlay(db, playID)
	if err != nil {
		slog.Error("Failed to get play", "error", err, "playID", playID)
		http.Error(w, "Failed to get play", http.StatusInternalServerError)
		return
	}

	if play == nil {
		http.Error(w, "Play not found", http.StatusNotFound)
		return
	}

	// Delete the play (this will also reorder remaining play numbers)
	if err := playrepo.DeletePlay(db, playID); err != nil {
		slog.Error("Failed to delete play", "error", err, "playID", playID)
		http.Error(w, "Failed to delete play", http.StatusInternalServerError)
		return
	}

	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d/games/%d/points/%d", matchID, setID, gameID, pointID))
	w.WriteHeader(http.StatusOK)
}

func GetByPoint(w http.ResponseWriter, r *http.Request) {
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

	plays, err := playrepo.GetPlaysByPoint(db, pointID)
	if err != nil {
		slog.Error("Failed to get plays", "error", err, "pointID", pointID)
		http.Error(w, "Failed to get plays", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plays)
	slog.Info("Handled", "method", r.Method, "path", r.URL.Path)
}

func Update(w http.ResponseWriter, r *http.Request) {
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

	playID := playshared.GetPlayID(w, r)
	if playID == 0 {
		slog.Error("Invalid play ID", "playID", playID)
		http.Error(w, "Invalid play ID", http.StatusBadRequest)
		return
	}

	// Get existing play
	existingPlay, err := playrepo.GetPlay(db, playID)
	if err != nil {
		slog.Error("Failed to get play", "error", err, "playID", playID)
		http.Error(w, "Failed to get play", http.StatusInternalServerError)
		return
	}

	if existingPlay == nil {
		http.Error(w, "Play not found", http.StatusNotFound)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Parse player ID
	playerIDStr := r.FormValue("player_id")
	if playerIDStr == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		slog.Error("Invalid player ID", "error", err, "playerID", playerIDStr)
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	// Parse ball positions
	ballXStr := r.FormValue("ball_position_x")
	ballYStr := r.FormValue("ball_position_y")
	ballX, err := strconv.Atoi(ballXStr)
	if err != nil || ballX < 0 || ballX > 6 {
		slog.Error("Invalid ball position X", "error", err, "ballX", ballXStr)
		http.Error(w, "Invalid ball position X", http.StatusBadRequest)
		return
	}
	ballY, err := strconv.Atoi(ballYStr)
	if err != nil || ballY < 0 || ballY > 11 {
		slog.Error("Invalid ball position Y", "error", err, "ballY", ballYStr)
		http.Error(w, "Invalid ball position Y", http.StatusBadRequest)
		return
	}

	// Get form values for the optional fields
	resultType := r.FormValue("result_type")
	handSide := r.FormValue("hand_side")
	contactType := r.FormValue("contact_type")
	shotEffect := r.FormValue("shot_effect")

	// Update the play
	updatedPlay := *existingPlay
	updatedPlay.PlayerID = sql.NullInt64{Int64: playerID, Valid: true}
	updatedPlay.BallPositionX = ballX
	updatedPlay.BallPositionY = ballY
	
	if resultType != "" && resultType != "Return" {
		updatedPlay.ResultType = sql.NullString{String: resultType, Valid: true}
	} else {
		updatedPlay.ResultType = sql.NullString{Valid: false}
	}
	
	if handSide != "" {
		updatedPlay.HandSide = sql.NullString{String: handSide, Valid: true}
	} else {
		updatedPlay.HandSide = sql.NullString{Valid: false}
	}
	
	if contactType != "" {
		updatedPlay.ContactType = sql.NullString{String: contactType, Valid: true}
	} else {
		updatedPlay.ContactType = sql.NullString{Valid: false}
	}
	
	if shotEffect != "" {
		updatedPlay.ShotEffect = sql.NullString{String: shotEffect, Valid: true}
	} else {
		updatedPlay.ShotEffect = sql.NullString{Valid: false}
	}

	// Save to database
	if err := playrepo.UpdatePlay(db, &updatedPlay); err != nil {
		slog.Error("Failed to update play", "error", err, "playID", playID)
		http.Error(w, "Failed to update play", http.StatusInternalServerError)
		return
	}

	// Check if this play ends the point (result_type is not null and not "Return")
	pointEnded := updatedPlay.ResultType.Valid && updatedPlay.ResultType.String != ""
	pointWasEnded := existingPlay.ResultType.Valid && existingPlay.ResultType.String != ""

	// Handle the case where point was previously ended but now continues
	if pointWasEnded && !pointEnded {
		slog.Warn("Point was previously ended but now continues", "pointID", pointID, "playID", playID, 
			"previousResult", existingPlay.ResultType.String)
		// Note: This could potentially affect subsequent points that were auto-created
		// For now, we allow this but log a warning. Future enhancement: validate consistency
	}
	
	if pointEnded {
		// Delete any subsequent plays in this point
		if err := playrepo.DeleteSubsequentPlays(db, pointID, updatedPlay.PlayNumber); err != nil {
			slog.Error("Failed to delete subsequent plays", "error", err, "pointID", pointID, "playNumber", updatedPlay.PlayNumber)
			http.Error(w, "Failed to cleanup plays", http.StatusInternalServerError)
			return
		}

		slog.Info("Point ended", "pointID", pointID, "finalResult", updatedPlay.ResultType.String)
		slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

		// Redirect back to the game (let user decide if game continues)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d/games/%d", matchID, setID, gameID))
		w.WriteHeader(http.StatusOK)
	} else {
		// Point continues - create next play in the same point
		allPlays, err := playrepo.GetPlaysByPoint(db, pointID)
		if err != nil {
			slog.Error("Failed to get plays for next play creation", "error", err, "pointID", pointID)
			http.Error(w, "Failed to get plays", http.StatusInternalServerError)
			return
		}

		nextPlay := playmodel.Play{
			PointID:       pointID,
			PlayNumber:    len(allPlays) + 1,
			PlayerID:      sql.NullInt64{Valid: false},
			BallPositionX: 0,
			BallPositionY: 0,
			ResultType:    sql.NullString{Valid: false},
			HandSide:      sql.NullString{Valid: false},
			ContactType:   sql.NullString{Valid: false},
			ShotEffect:    sql.NullString{Valid: false},
		}

		if err := playrepo.CreatePlay(db, &nextPlay); err != nil {
			slog.Error("Failed to create next play", "error", err, "pointID", pointID)
			http.Error(w, "Failed to create next play", http.StatusInternalServerError)
			return
		}

		slog.Info("Point continues, created next play", "pointID", pointID, "nextPlayID", nextPlay.ID, "playNumber", nextPlay.PlayNumber)
		slog.Info("Handled", "method", r.Method, "path", r.URL.Path)

		// Redirect to the new play
		w.Header().Set("HX-Redirect", fmt.Sprintf("/matches/%d/sets/%d/games/%d/points/%d/plays/%d", matchID, setID, gameID, pointID, nextPlay.ID))
		w.WriteHeader(http.StatusOK)
	}
}