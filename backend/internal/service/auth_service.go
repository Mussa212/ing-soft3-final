package service

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"

	servicedto "vesuvio/internal/dto/service"
)

// UserClient abstracts user persistence.
type UserClient interface {
	CreateUser(ctx context.Context, params servicedto.CreateUserParams) (*servicedto.User, error)
	GetUserByEmail(ctx context.Context, email string) (*servicedto.UserWithPassword, error)
	GetUserByID(ctx context.Context, id uint) (*servicedto.User, error)
}

type AuthService struct {
	userClient UserClient
}

func NewAuthService(userClient UserClient) *AuthService {
	return &AuthService{userClient: userClient}
}

func (s *AuthService) Register(ctx context.Context, input servicedto.RegisterUserInput) (*servicedto.RegisterUserOutput, error) {
	return nil, nil
}

func (s *AuthService) Login(ctx context.Context, input servicedto.LoginUserInput) (*servicedto.LoginUserOutput, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || strings.TrimSpace(input.Password) == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &servicedto.LoginUserOutput{User: user.User}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uint) (*servicedto.User, error) {
	user, err := s.userClient.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
