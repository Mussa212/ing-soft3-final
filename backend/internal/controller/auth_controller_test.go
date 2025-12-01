package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	controllerdto "vesuvio/internal/dto/controller"
	servicedto "vesuvio/internal/dto/service"
	"vesuvio/internal/service"
)

func TestAuthController_RegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authSvc := service.NewAuthService(newControllerFakeUserClient())
	ctl := NewAuthController(authSvc)

	// Register
	registerBody := controllerdto.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctl.Register(newTestContext(req, w))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var registerResp controllerdto.RegisterResponse
	if err := json.Unmarshal(w.Body.Bytes(), &registerResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if registerResp.Email != registerBody.Email || registerResp.ID == 0 {
		t.Fatalf("unexpected register response: %+v", registerResp)
	}

	// Login success
	loginBody := controllerdto.LoginRequest{
		Email:    "alice@example.com",
		Password: "secret",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	ctl.Login(newTestContext(req, w))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var loginResp controllerdto.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if loginResp.ID != registerResp.ID {
		t.Fatalf("expected login id %d, got %d", registerResp.ID, loginResp.ID)
	}
}

func TestAuthController_DuplicateEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userClient := newControllerFakeUserClient()
	authSvc := service.NewAuthService(userClient)
	ctl := NewAuthController(authSvc)

	// Seed existing
	_, _ = userClient.CreateUser(context.Background(), servicedto.CreateUserParams{
		Name:         "Existing",
		Email:        "dupe@example.com",
		PasswordHash: "hash",
	})

	body, _ := json.Marshal(controllerdto.RegisterRequest{
		Name:     "New",
		Email:    "dupe@example.com",
		Password: "secret",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctl.Register(newTestContext(req, w))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthController_LoginInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userClient := newControllerFakeUserClient()
	authSvc := service.NewAuthService(userClient)
	ctl := NewAuthController(authSvc)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	_, _ = userClient.CreateUser(context.Background(), servicedto.CreateUserParams{
		Name:         "Bob",
		Email:        "bob@example.com",
		PasswordHash: string(hash),
	})

	body, _ := json.Marshal(controllerdto.LoginRequest{
		Email:    "bob@example.com",
		Password: "wrong",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctl.Login(newTestContext(req, w))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// Helpers and fakes for controller tests.

func newTestContext(req *http.Request, w http.ResponseWriter) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}

type controllerFakeUserClient struct {
	users  map[uint]servicedto.UserWithPassword
	nextID uint
}

func newControllerFakeUserClient() *controllerFakeUserClient {
	return &controllerFakeUserClient{
		users:  make(map[uint]servicedto.UserWithPassword),
		nextID: 1,
	}
}

func (f *controllerFakeUserClient) CreateUser(ctx context.Context, params servicedto.CreateUserParams) (*servicedto.User, error) {
	for _, u := range f.users {
		if u.Email == params.Email {
			return nil, service.ErrEmailAlreadyExists
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

func (f *controllerFakeUserClient) GetUserByEmail(ctx context.Context, email string) (*servicedto.UserWithPassword, error) {
	for _, u := range f.users {
		if u.Email == email {
			copy := u
			return &copy, nil
		}
	}
	return nil, nil
}

func (f *controllerFakeUserClient) GetUserByID(ctx context.Context, id uint) (*servicedto.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, nil
	}
	copy := u.User
	return &copy, nil
}
