package model

import "time"

// ReservationModel represents a booking in the system.
type ReservationModel struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	User      UserModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Date      time.Time `gorm:"type:date;not null"`
	Time      string    `gorm:"size:5;not null"` // HH:MM
	People    int       `gorm:"not null"`
	Comment   *string   `gorm:"type:text"`
	Status    string    `gorm:"size:20;not null;default:pending"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
