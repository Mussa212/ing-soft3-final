package client

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"vesuvio/internal/model"
)

// Migrate ensures database tables exist for all models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.UserModel{}, &model.ReservationModel{})
}

// SeedUser describes a user to insert.
type SeedUser struct {
	Name     string
	Email    string
	Password string
	IsAdmin  bool
}

// SeedUsers inserts users if they don't already exist (by email).
func SeedUsers(ctx context.Context, db *gorm.DB, users []SeedUser) error {
	for _, u := range users {
		email := strings.TrimSpace(strings.ToLower(u.Email))
		if email == "" || strings.TrimSpace(u.Password) == "" {
			return errors.New("seed user missing email or password")
		}

		var existing model.UserModel
		err := db.WithContext(ctx).Where("email = ?", email).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existing.ID != 0 {
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newUser := model.UserModel{
			Name:         strings.TrimSpace(u.Name),
			Email:        email,
			PasswordHash: string(hash),
			IsAdmin:      u.IsAdmin,
		}
		if err := db.WithContext(ctx).Create(&newUser).Error; err != nil {
			return err
		}
	}
	return nil
}
