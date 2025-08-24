package database

import (
	"ct-padel-s/src/infrastructure/env"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

var instance *DB

func Initialize() (*DB, error) {
	if instance != nil {
		return instance, nil
	}

	db, err := sql.Open("postgres", env.DatabaseURL)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	instance = &DB{db}
	slog.Info("Database connection established successfully")

	return instance, nil
}

func GetDB() *DB {
	if instance == nil {
		panic("database not initialized - call Initialize() first")
	}
	return instance
}

func (db *DB) Close() error {
	if db.DB != nil {
		slog.Info("Closing database connection")
		return db.DB.Close()
	}
	return nil
}

func (db *DB) IsHealthy() bool {
	if db.DB == nil {
		return false
	}
	return db.Ping() == nil
}
