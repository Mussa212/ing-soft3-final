package client

import (
	"testing"

	"gorm.io/driver/sqlite"
)

// Ensures NewDBWithDialector runs migrations and returns a working DB.
func TestNewDBWithDialectorMigrations(t *testing.T) {
	db, err := NewDBWithDialector(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		t.Fatalf("unexpected error opening sqlite db: %v", err)
	}

	if !db.Migrator().HasTable("user_models") || !db.Migrator().HasTable("reservation_models") {
		t.Fatalf("expected migrations to create user_models and reservation_models tables")
	}
}

// NewDB should surface an error for an invalid DSN (covers the wrapper).
func TestNewDBInvalidDSN(t *testing.T) {
	if _, err := NewDB("not a dsn"); err == nil {
		t.Fatalf("expected error for invalid DSN")
	}
}
