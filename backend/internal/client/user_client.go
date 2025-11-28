package client

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"vesuvio/internal/dto/service"
	"vesuvio/internal/model"
)

type GormUserClient struct {
	db *gorm.DB
}

func NewUserClient(db *gorm.DB) *GormUserClient {
	return &GormUserClient{db: db}
}

func (c *GormUserClient) CreateUser(ctx context.Context, params servicedto.CreateUserParams) (*servicedto.User, error) {
	user := model.UserModel{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		IsAdmin:      params.IsAdmin,
	}

	if err := c.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return toServiceUser(&user), nil
}

func (c *GormUserClient) GetUserByEmail(ctx context.Context, email string) (*servicedto.UserWithPassword, error) {
	var user model.UserModel
	err := c.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &servicedto.UserWithPassword{
		User:         *toServiceUser(&user),
		PasswordHash: user.PasswordHash,
	}, nil
}

func (c *GormUserClient) GetUserByID(ctx context.Context, id uint) (*servicedto.User, error) {
	var user model.UserModel
	err := c.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toServiceUser(&user), nil
}

func toServiceUser(u *model.UserModel) *servicedto.User {
	return &servicedto.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
