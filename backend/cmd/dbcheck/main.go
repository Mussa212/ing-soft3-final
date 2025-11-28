package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

// Small helper to verify DB connectivity using pgx directly.
func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("query failed: %v", err)
	}

	log.Println("Connected to:", version)
}
