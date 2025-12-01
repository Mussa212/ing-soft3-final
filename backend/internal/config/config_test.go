package config

import "testing"

// Ensures Load returns defaults when environment variables are empty or unset.
func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_DSN", "")

	cfg := Load()
	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.DatabaseDSN != "postgres://postgres:postgres@localhost:5432/vesuvio?sslmode=disable" {
		t.Fatalf("unexpected default dsn: %s", cfg.DatabaseDSN)
	}
}

// Ensures Load respects provided environment variables.
func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("DATABASE_DSN", "postgres://user:pass@db.example.com:5432/otherdb?sslmode=disable")

	cfg := Load()
	if cfg.Port != "9000" {
		t.Fatalf("expected port 9000, got %s", cfg.Port)
	}
	if cfg.DatabaseDSN != "postgres://user:pass@db.example.com:5432/otherdb?sslmode=disable" {
		t.Fatalf("unexpected dsn: %s", cfg.DatabaseDSN)
	}
}
