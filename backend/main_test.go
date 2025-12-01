package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vesuvio/internal/client"
)

func TestRedactDSN(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
			name: "masks password",
			dsn:  "postgres://user:secret@localhost:5432/db?sslmode=disable",
			want: "postgres://user:%2A%2A%2A%2A@localhost:5432/db?sslmode=disable",
		},
		{
			name: "returns original when no credentials",
			dsn:  "postgres://localhost:5432/db",
			want: "postgres://localhost:5432/db",
		},
		{
			name: "returns original when parse fails",
			dsn:  "not a valid dsn %",
			want: "not a valid dsn %",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactDSN(tt.dsn); got != tt.want {
				t.Fatalf("redactDSN() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMainWithStubbedDependencies(t *testing.T) {
	origOpenDB := openDB
	origStartHTTP := startHTTP
	defer func() {
		openDB = origOpenDB
		startHTTP = origStartHTTP
	}()

	openDB = func(dsn string) (*gorm.DB, error) {
		return client.NewDBWithDialector(sqlite.Open("file:main_test?mode=memory&cache=shared"))
	}

	var started bool
	startHTTP = func(r *gin.Engine, port string) error {
		started = true
		if port == "" {
			t.Fatalf("expected port to be set")
		}
		return nil
	}

	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_DSN", "ignored")

	main()

	if !started {
		t.Fatalf("expected startHTTP to be called")
	}
}
