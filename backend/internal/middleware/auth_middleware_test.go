package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	servicedto "vesuvio/internal/dto/service"
	"vesuvio/internal/service"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authSvc := service.NewAuthService(newMiddlewareFakeUserClient())

	r := gin.New()
	r.Use(AuthMiddleware(authSvc))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authSvc := service.NewAuthService(newMiddlewareFakeUserClient())

	r := gin.New()
	r.Use(AuthMiddleware(authSvc))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", "99")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAdminOnly_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userClient := newMiddlewareFakeUserClient()
	// Seed non-admin user with id 1.
	_, _ = userClient.CreateUser(context.Background(), servicedto.CreateUserParams{
		Name:         "User",
		Email:        "user@example.com",
		PasswordHash: "hash",
		IsAdmin:      false,
	})
	authSvc := service.NewAuthService(userClient)

	r := gin.New()
	r.Use(AuthMiddleware(authSvc), AdminOnly())
	r.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// fake user client for middleware tests.
type middlewareFakeUserClient struct {
	users map[uint]servicedto.UserWithPassword
}

func newMiddlewareFakeUserClient() *middlewareFakeUserClient {
	return &middlewareFakeUserClient{users: make(map[uint]servicedto.UserWithPassword)}
}

func (f *middlewareFakeUserClient) CreateUser(ctx context.Context, params servicedto.CreateUserParams) (*servicedto.User, error) {
	id := uint(len(f.users) + 1)
	user := servicedto.UserWithPassword{
		User: servicedto.User{
			ID:      id,
			Name:    params.Name,
			Email:   params.Email,
			IsAdmin: params.IsAdmin,
		},
		PasswordHash: params.PasswordHash,
	}
	f.users[id] = user
	return &user.User, nil
}

func (f *middlewareFakeUserClient) GetUserByEmail(ctx context.Context, email string) (*servicedto.UserWithPassword, error) {
	for _, u := range f.users {
		if u.Email == email {
			copy := u
			return &copy, nil
		}
	}
	return nil, nil
}

func (f *middlewareFakeUserClient) GetUserByID(ctx context.Context, id uint) (*servicedto.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, nil
	}
	copy := u.User
	return &copy, nil
}
