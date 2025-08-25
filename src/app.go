package src

import (
	"ct-padel-s/src/features/home"
	"ct-padel-s/src/features/padel/game"
	"ct-padel-s/src/features/padel/match"
	"ct-padel-s/src/features/padel/play"
	"ct-padel-s/src/features/padel/point"
	"ct-padel-s/src/features/padel/set"
	"ct-padel-s/src/infrastructure/database"
	"ct-padel-s/src/infrastructure/fileserver"
	"log"
	"log/slog"
	"net/http"
)

func App() {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Routes
	mux := http.NewServeMux()

	// Hypermedia routes (HTML)
	mux.HandleFunc("GET /matches", match.GetAll)
	mux.HandleFunc("POST /matches", match.Create)
	mux.HandleFunc("GET /matches/{matchID}", match.Get)
	mux.HandleFunc("DELETE /matches/{matchID}", match.Delete)

	mux.HandleFunc("POST /matches/{matchID}/sets", set.Create)
	mux.HandleFunc("GET /matches/{matchID}/sets/{setID}", set.Get)
	mux.HandleFunc("DELETE /matches/{matchID}/sets/{setID}", set.Delete)

	mux.HandleFunc("POST /matches/{matchID}/sets/{setID}/games", game.Create)
	mux.HandleFunc("GET /matches/{matchID}/sets/{setID}/games/{gameID}", game.Get)
	mux.HandleFunc("DELETE /matches/{matchID}/sets/{setID}/games/{gameID}", game.Delete)

	mux.HandleFunc("POST /matches/{matchID}/sets/{setID}/games/{gameID}/points", point.Create)
	mux.HandleFunc("GET /matches/{matchID}/sets/{setID}/games/{gameID}/points/{pointID}", point.Get)
	mux.HandleFunc("DELETE /matches/{matchID}/sets/{setID}/games/{gameID}/points/{pointID}", point.Delete)

	mux.HandleFunc("POST /matches/{matchID}/sets/{setID}/games/{gameID}/points/{pointID}/plays", play.Create)
	mux.HandleFunc("GET /matches/{matchID}/sets/{setID}/games/{gameID}/points/{pointID}/plays/{playID}", play.Get)
	mux.HandleFunc("PATCH /matches/{matchID}/sets/{setID}/games/{gameID}/points/{pointID}/plays/{playID}", play.Patch)
	mux.HandleFunc("DELETE /matches/{matchID}/sets/{setID}/games/{gameID}/points/{pointID}/plays/{playID}", play.Delete)

	// Home page
	mux.HandleFunc("/", home.Handler)

	// Handle Chrome DevTools well-known endpoint to keep logs clean
	mux.HandleFunc("/.well-known/appspecific/com.chrome.devtools.json", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Static files with ETag caching
	cachedFS := fileserver.NewCachedFileServer("build/static/")
	mux.Handle("/static/", http.StripPrefix("/static/", cachedFS))

	slog.Info("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
