package client

import (
	"context"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
)

func TestMigrateAndSeedUsers(t *testing.T) {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := NewDBWithDialector(sqlite.Open(dsn))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	users := []SeedUser{
		{Name: "Alice", Email: "alice@example.com", Password: "pw1", IsAdmin: true},
		{Name: "Bob", Email: "bob@example.com", Password: "pw2"},
	}
	if err := SeedUsers(context.Background(), db, users); err != nil {
		t.Fatalf("seed users: %v", err)
	}

	client := NewUserClient(db)
	found, err := client.GetUserByEmail(context.Background(), "alice@example.com")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if found == nil || !found.IsAdmin {
		t.Fatalf("expected seeded admin user, got %+v", found)
	}

	// Running seed again should be idempotent (no duplicates).
	if err := SeedUsers(context.Background(), db, users); err != nil {
		t.Fatalf("seed users second time: %v", err)
	}
}

func TestSeedUsersValidation(t *testing.T) {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := NewDBWithDialector(sqlite.Open(dsn))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	err = SeedUsers(context.Background(), db, []SeedUser{
		{Name: "Bad", Email: "   ", Password: "pw"},
	})
	if err == nil {
		t.Fatalf("expected error for missing email")
	}

	err = SeedUsers(context.Background(), db, []SeedUser{
		{Name: "Bad", Email: "bad@example.com", Password: "   "},
	})
	if err == nil {
		t.Fatalf("expected error for missing password")
	}
}
