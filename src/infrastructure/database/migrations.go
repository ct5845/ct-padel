package database

import (
	"embed"
	"fmt"
	"log/slog"
	"sort"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

var migrations = []Migration{}

func RegisterMigration(version int, name, up, down string) {
	migrations = append(migrations, Migration{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

func RunMigrations(db *DB) error {
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, migration := range migrations {
		if _, applied := appliedMigrations[migration.Version]; applied {
			slog.Debug("Migration already applied", "version", migration.Version, "name", migration.Name)
			continue
		}

		slog.Info("Running migration", "version", migration.Version, "name", migration.Name)

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
		}

		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", migration.Version, migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		slog.Info("Migration completed", "version", migration.Version, "name", migration.Name)
	}

	return nil
}

func createMigrationsTable(db *DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

func getAppliedMigrations(db *DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func init() {
	v001Up, _ := migrationFiles.ReadFile("migrations/001_up.sql")
	v001Down, _ := migrationFiles.ReadFile("migrations/001_down.sql")

	RegisterMigration(1, "create_padel_tables", string(v001Up), string(v001Down))
}
