package service

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	servicedto "vesuvio/internal/dto/service"
)

func TestRegisterDuplicateEmail(t *testing.T) {
	userClient := newFakeUserClient()
	ctx := context.Background()

	_, err := userClient.CreateUser(ctx, servicedto.CreateUserParams{
		Name:         "Existing",
		Email:        "test@example.com",
		PasswordHash: "hash",
		IsAdmin:      false,
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	svc := NewAuthService(userClient)
	_, err = svc.Register(ctx, servicedto.RegisterUserInput{
		Name:     "New",
		Email:    "test@example.com",
		Password: "secret",
	})
	if err != ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	userClient := newFakeUserClient()
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	_, err := userClient.CreateUser(ctx, servicedto.CreateUserParams{
		Name:         "User",
		Email:        "login@example.com",
		PasswordHash: string(hash),
		IsAdmin:      false,
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	svc := NewAuthService(userClient)
	_, err = svc.Login(ctx, servicedto.LoginUserInput{
		Email:    "login@example.com",
		Password: "wrong",
	})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestRegisterSuccessAndLoginSuccess(t *testing.T) {
	userClient := newFakeUserClient()
	ctx := context.Background()

	svc := NewAuthService(userClient)
	registerOut, err := svc.Register(ctx, servicedto.RegisterUserInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}
	if registerOut.User.Email != "alice@example.com" || registerOut.User.ID == 0 {
		t.Fatalf("unexpected register output: %+v", registerOut.User)
	}

	loginOut, err := svc.Login(ctx, servicedto.LoginUserInput{
		Email:    "alice@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
	if loginOut.User.ID != registerOut.User.ID {
		t.Fatalf("expected same user id, got %d vs %d", loginOut.User.ID, registerOut.User.ID)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	userClient := newFakeUserClient()
	svc := NewAuthService(userClient)
	ctx := context.Background()

	_, err := svc.GetUserByID(ctx, 123)
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

// fakeUserClient is a simple in-memory implementation for tests.
type fakeUserClient struct {
	users  map[uint]servicedto.UserWithPassword
	nextID uint
}

func newFakeUserClient() *fakeUserClient {
	return &fakeUserClient{
		users:  make(map[uint]servicedto.UserWithPassword),
		nextID: 1,
	}
}

func (f *fakeUserClient) CreateUser(ctx context.Context, params servicedto.CreateUserParams) (*servicedto.User, error) {
	for _, u := range f.users {
		if u.Email == params.Email {
			return nil, ErrEmailAlreadyExists
		}
	}
	id := f.nextID
	f.nextID++

	now := time.Now()
	user := servicedto.UserWithPassword{
		User: servicedto.User{
			ID:        id,
			Name:      params.Name,
			Email:     params.Email,
			IsAdmin:   params.IsAdmin,
			CreatedAt: now,
			UpdatedAt: now,
		},
		PasswordHash: params.PasswordHash,
	}
	f.users[id] = user
	return &user.User, nil
}

func (f *fakeUserClient) GetUserByEmail(ctx context.Context, email string) (*servicedto.UserWithPassword, error) {
	for _, u := range f.users {
		if u.Email == email {
			copy := u
			return &copy, nil
		}
	}
	return nil, nil
}

func (f *fakeUserClient) GetUserByID(ctx context.Context, id uint) (*servicedto.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, nil
	}
	copy := u.User
	return &copy, nil
}
