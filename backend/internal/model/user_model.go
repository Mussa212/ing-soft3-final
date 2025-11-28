package model

import "time"

// UserModel represents the persisted user.
type UserModel struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:255;not null"`
	Email        string `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash string `gorm:"not null"`
	IsAdmin      bool   `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
