package client

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"vesuvio/internal/dto/service"
	"vesuvio/internal/model"
)

type GormReservationClient struct {
	db *gorm.DB
}

func NewReservationClient(db *gorm.DB) *GormReservationClient {
	return &GormReservationClient{db: db}
}

func (c *GormReservationClient) CreateReservation(ctx context.Context, params servicedto.CreateReservationParams) (*servicedto.Reservation, error) {
	res := model.ReservationModel{
		UserID:  params.UserID,
		Date:    params.Date,
		Time:    params.Time,
		People:  params.People,
		Comment: params.Comment,
		Status:  params.Status,
	}

	if err := c.db.WithContext(ctx).Create(&res).Error; err != nil {
		return nil, err
	}
	return toServiceReservation(&res, nil), nil
}

func (c *GormReservationClient) ListReservationsByUser(ctx context.Context, userID uint, status *string) ([]servicedto.Reservation, error) {
	var models []model.ReservationModel
	query := c.db.WithContext(ctx).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Order("date, time").Find(&models).Error; err != nil {
		return nil, err
	}
	return mapReservations(models, nil), nil
}

func (c *GormReservationClient) GetReservationByID(ctx context.Context, id uint) (*servicedto.Reservation, error) {
	var res model.ReservationModel
	err := c.db.WithContext(ctx).First(&res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toServiceReservation(&res, nil), nil
}

func (c *GormReservationClient) UpdateReservationStatus(ctx context.Context, id uint, status string) (*servicedto.Reservation, error) {
	var res model.ReservationModel
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&res, id).Error; err != nil {
			return err
		}
		res.Status = status
		return tx.Save(&res).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toServiceReservation(&res, nil), nil
}

func (c *GormReservationClient) ListReservationsByDate(ctx context.Context, date time.Time, status *string) ([]servicedto.Reservation, error) {
	var models []model.ReservationModel
	query := c.db.WithContext(ctx).Preload("User").Where("date = ?", date)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Order("time").Find(&models).Error; err != nil {
		return nil, err
	}
	return mapReservations(models, func(m model.ReservationModel) *servicedto.User {
		return toServiceUser(&m.User)
	}), nil
}

func mapReservations(models []model.ReservationModel, userMapper func(model.ReservationModel) *servicedto.User) []servicedto.Reservation {
	reservations := make([]servicedto.Reservation, 0, len(models))
	for _, m := range models {
		var user *servicedto.User
		if userMapper != nil {
			user = userMapper(m)
		}
		reservations = append(reservations, *toServiceReservation(&m, user))
	}
	return reservations
}

func toServiceReservation(m *model.ReservationModel, user *servicedto.User) *servicedto.Reservation {
	return &servicedto.Reservation{
		ID:        m.ID,
		UserID:    m.UserID,
		User:      user,
		Date:      m.Date,
		Time:      m.Time,
		People:    m.People,
		Comment:   m.Comment,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
