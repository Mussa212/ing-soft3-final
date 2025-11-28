package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds app level configuration.
type Config struct {
	Port        string
	DatabaseDSN string
}

// Load returns configuration using environment variables with sane defaults.
func Load() Config {
	// Try loading .env from current directory or parent directories
	if err := godotenv.Load(); err != nil {
		// If not found in current dir, try going up two levels (useful for cmd/seed)
		_ = godotenv.Load("../../.env")
	}

	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/vesuvio?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
