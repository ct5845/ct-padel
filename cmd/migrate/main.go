package main

import (
	"ct-padel-s/src/infrastructure/database"
	_ "ct-padel-s/src/infrastructure/logging" // Import for colored logging init
	"log/slog"
)

func main() {
	db, err := database.Initialize()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		panic(err)
	}

	slog.Info("Migrations completed successfully")
}
