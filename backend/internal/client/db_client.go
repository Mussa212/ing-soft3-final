package client

import (
	"vesuvio/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB opens a GORM connection and runs migrations.
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.UserModel{}, &model.ReservationModel{}); err != nil {
		return nil, err
	}

	return db, nil
}
