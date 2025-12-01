package client

import (
	"context"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	servicedto "vesuvio/internal/dto/service"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := NewDBWithDialector(sqlite.Open(dsn))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db
}

func TestUserClient_CreateAndGet(t *testing.T) {
	db := newTestDB(t)
	client := NewUserClient(db)
	ctx := context.Background()

	created, err := client.CreateUser(ctx, servicedto.CreateUserParams{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hash",
		IsAdmin:      true,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if created.ID == 0 || created.Email != "alice@example.com" || !created.IsAdmin {
		t.Fatalf("unexpected created user: %+v", created)
	}

	byEmail, err := client.GetUserByEmail(ctx, "alice@example.com")
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}
	if byEmail == nil || byEmail.PasswordHash != "hash" {
		t.Fatalf("unexpected user by email: %+v", byEmail)
	}

	byID, err := client.GetUserByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if byID == nil || byID.Email != "alice@example.com" {
		t.Fatalf("unexpected user by id: %+v", byID)
	}
}

func TestUserClient_DuplicateEmailAndNotFound(t *testing.T) {
	db := newTestDB(t)
	client := NewUserClient(db)
	ctx := context.Background()

	_, err := client.CreateUser(ctx, servicedto.CreateUserParams{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hash",
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	_, err = client.CreateUser(ctx, servicedto.CreateUserParams{
		Name:         "Alice 2",
		Email:        "alice@example.com",
		PasswordHash: "hash2",
	})
	if err == nil {
		t.Fatalf("expected duplicate email error")
	}

	user, err := client.GetUserByEmail(ctx, "missing@example.com")
	if err != nil {
		t.Fatalf("unexpected error for missing email: %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil for missing email, got %+v", user)
	}

	userByID, err := client.GetUserByID(ctx, 999)
	if err != nil {
		t.Fatalf("unexpected error for missing id: %v", err)
	}
	if userByID != nil {
		t.Fatalf("expected nil for missing id, got %+v", userByID)
	}
}
