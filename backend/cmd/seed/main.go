package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"vesuvio/internal/client"
	"vesuvio/internal/config"
)

// Seed command to run migrations and insert default admin/user.
func main() {
	cfg := config.Load()

	log.Printf("Using DATABASE_DSN: %s", redactDSN(cfg.DatabaseDSN))
	db, err := client.NewDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := client.Migrate(db); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminName := getenvDefault("ADMIN_NAME", "Admin")

	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")
	userName := getenvDefault("USER_NAME", "User")

	if adminEmail == "" || adminPassword == "" || userEmail == "" || userPassword == "" {
		log.Fatalf("missing required env vars ADMIN_EMAIL/ADMIN_PASSWORD and USER_EMAIL/USER_PASSWORD")
	}

	ctx := context.Background()
	if err := client.SeedUsers(ctx, db, []client.SeedUser{
		{Name: adminName, Email: adminEmail, Password: adminPassword, IsAdmin: true},
		{Name: userName, Email: userEmail, Password: userPassword, IsAdmin: false},
	}); err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}

	log.Println("migration and seeding completed")
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// redactDSN masks the password in the DSN for logging.
func redactDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil || u.User == nil {
		return dsn
	}
	username := u.User.Username()
	u.User = url.UserPassword(username, "****")
	return u.String()
}
