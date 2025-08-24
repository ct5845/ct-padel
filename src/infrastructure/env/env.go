package env

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var DatabaseURL string

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func init() {
	slog.Debug("Environment variables initialized", "component", "env")
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load environment variables", "error", err)
		panic(err)
	}

	DatabaseURL = getEnv("DATABASE_URL", "")
}
