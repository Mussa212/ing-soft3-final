package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	controllerdto "vesuvio/internal/dto/controller"
	servicedto "vesuvio/internal/dto/service"
	"vesuvio/internal/middleware"
	"vesuvio/internal/service"
)

func TestReservationController_UserFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userClient := newControllerFakeUserClient()
	authSvc := service.NewAuthService(userClient)

	// Seed user
	user, _ := userClient.CreateUser(context.Background(), servicedto.CreateUserParams{
		Name:         "User",
		Email:        "user@example.com",
		PasswordHash: "hash",
	})

	resClient := newControllerFakeReservationClient()
	resSvc := service.NewReservationService(resClient)
	resCtl := NewReservationController(resSvc)

	router := gin.New()
	router.Use(middleware.AuthMiddleware(authSvc))
	router.POST("/reservations", resCtl.CreateReservation)
	router.GET("/my/reservations", resCtl.ListMyReservations)
	router.PATCH("/reservations/:id/cancel", resCtl.CancelReservation)

	// Create reservation
	payload := controllerdto.CreateReservationRequest{
		Date:   "2025-12-01",
		Time:   "20:30",
		People: 4,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", user.ID))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created controllerdto.ReservationResponse
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	// List my reservations
	req = httptest.NewRequest(http.MethodGet, "/my/reservations", nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", user.ID))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []controllerdto.ReservationResponse
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(list))
	}

	// Cancel reservation
	cancelURL := fmt.Sprintf("/reservations/%d/cancel", created.ID)
	req = httptest.NewRequest(http.MethodPatch, cancelURL, nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", user.ID))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var cancelled controllerdto.ReservationResponse
	_ = json.Unmarshal(w.Body.Bytes(), &cancelled)
	if cancelled.Status != servicedto.StatusCancelled {
		t.Fatalf("expected cancelled, got %s", cancelled.Status)
	}
}

func TestReservationController_ErrorBranches(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resClient := newControllerFakeReservationClient()
	resSvc := service.NewReservationService(resClient)
	resCtl := NewReservationController(resSvc)

	// Create reservation with invalid payload (missing people)
	body, _ := json.Marshal(controllerdto.CreateReservationRequest{
		Date: "2025-12-01",
		Time: "20:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := newTestContext(req, w)
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 1})
	resCtl.CreateReservation(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid input, got %d", w.Code)
	}

	// List with invalid status
	req = httptest.NewRequest(http.MethodGet, "/my/reservations?status=weird", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 1})
	resCtl.ListMyReservations(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid status, got %d", w.Code)
	}

	// Cancel with invalid ID
	req = httptest.NewRequest(http.MethodPatch, "/reservations/abc/cancel", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 1})
	resCtl.CancelReservation(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id, got %d", w.Code)
	}

	// Seed reservation belonging to another user to trigger forbidden
	otherRes, _ := resClient.CreateReservation(context.Background(), servicedto.CreateReservationParams{
		UserID:  2,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "19:00",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})
	req = httptest.NewRequest(http.MethodPatch, "/reservations/forbidden/cancel", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: fmt.Sprintf("%d", otherRes.ID)}}
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 1})
	resCtl.CancelReservation(c)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for forbidden cancel, got %d", w.Code)
	}

	// Not found case
	req = httptest.NewRequest(http.MethodPatch, "/reservations/missing/cancel", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 1})
	resCtl.CancelReservation(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing reservation, got %d", w.Code)
	}
}

func TestReservationController_AdminFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userClient := newControllerFakeUserClient()
	authSvc := service.NewAuthService(userClient)
	admin, _ := userClient.CreateUser(context.Background(), servicedto.CreateUserParams{
		Name:         "Admin",
		Email:        "admin@example.com",
		PasswordHash: "hash",
		IsAdmin:      true,
	})
	user, _ := userClient.CreateUser(context.Background(), servicedto.CreateUserParams{
		Name:         "User",
		Email:        "user@example.com",
		PasswordHash: "hash",
	})

	resClient := newControllerFakeReservationClient()
	resSvc := service.NewReservationService(resClient)
	adminCtl := NewAdminController(resSvc)

	// Seed reservation
	_, _ = resClient.CreateReservation(context.Background(), servicedto.CreateReservationParams{
		UserID:  user.ID,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "20:00",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})

	router := gin.New()
	router.Use(middleware.AuthMiddleware(authSvc), middleware.AdminOnly())
	router.GET("/admin/reservations", adminCtl.ListReservations)
	router.PATCH("/admin/reservations/:id/confirm", adminCtl.ConfirmReservation)
	router.PATCH("/admin/reservations/:id/cancel", adminCtl.CancelReservation)

	// List reservations
	req := httptest.NewRequest(http.MethodGet, "/admin/reservations?date=2025-12-01", nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []controllerdto.AdminReservationResponse
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(list))
	}

	// Confirm reservation
	confirmURL := fmt.Sprintf("/admin/reservations/%d/confirm", list[0].ID)
	req = httptest.NewRequest(http.MethodPatch, confirmURL, nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on confirm, got %d", w.Code)
	}

	// Cancel reservation
	cancelURL := fmt.Sprintf("/admin/reservations/%d/cancel", list[0].ID)
	req = httptest.NewRequest(http.MethodPatch, cancelURL, nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on cancel, got %d", w.Code)
	}
}

func TestAdminController_ErrorBranches(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resSvc := service.NewReservationService(newControllerFakeReservationClient())
	adminCtl := NewAdminController(resSvc)

	// List missing date
	req := httptest.NewRequest(http.MethodGet, "/admin/reservations", nil)
	w := httptest.NewRecorder()
	c := newTestContext(req, w)
	adminCtl.ListReservations(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing date, got %d", w.Code)
	}

	// List invalid status
	req = httptest.NewRequest(http.MethodGet, "/admin/reservations?date=2025-12-01&status=weird", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	adminCtl.ListReservations(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid status, got %d", w.Code)
	}

	// Confirm invalid id
	req = httptest.NewRequest(http.MethodPatch, "/admin/reservations/abc/confirm", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 1, IsAdmin: true})
	adminCtl.ConfirmReservation(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id, got %d", w.Code)
	}

	// Confirm unauthorized (non-admin)
	req = httptest.NewRequest(http.MethodPatch, "/admin/reservations/1/confirm", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 2, IsAdmin: false})
	adminCtl.ConfirmReservation(c)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for unauthorized admin, got %d", w.Code)
	}

	// Cancel not found
	req = httptest.NewRequest(http.MethodPatch, "/admin/reservations/999/cancel", nil)
	w = httptest.NewRecorder()
	c = newTestContext(req, w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}
	c.Set(middleware.ContextUserKey, servicedto.User{ID: 1, IsAdmin: true})
	adminCtl.CancelReservation(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing reservation, got %d", w.Code)
	}
}

// Fake reservation client for controller tests.
type controllerFakeReservationClient struct {
	reservations map[uint]servicedto.Reservation
	nextID       uint
}

func newControllerFakeReservationClient() *controllerFakeReservationClient {
	return &controllerFakeReservationClient{
		reservations: make(map[uint]servicedto.Reservation),
		nextID:       1,
	}
}

func (f *controllerFakeReservationClient) CreateReservation(ctx context.Context, params servicedto.CreateReservationParams) (*servicedto.Reservation, error) {
	id := f.nextID
	f.nextID++
	now := time.Now()
	res := servicedto.Reservation{
		ID:     id,
		UserID: params.UserID,
		User: &servicedto.User{
			ID:    params.UserID,
			Name:  fmt.Sprintf("User%d", params.UserID),
			Email: fmt.Sprintf("user%d@example.com", params.UserID),
		},
		Date:      params.Date,
		Time:      params.Time,
		People:    params.People,
		Comment:   params.Comment,
		Status:    params.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	f.reservations[id] = res
	return &res, nil
}

func (f *controllerFakeReservationClient) ListReservationsByUser(ctx context.Context, userID uint, status *string) ([]servicedto.Reservation, error) {
	var list []servicedto.Reservation
	for _, r := range f.reservations {
		if r.UserID != userID {
			continue
		}
		if status != nil && r.Status != *status {
			continue
		}
		list = append(list, r)
	}
	return list, nil
}

func (f *controllerFakeReservationClient) GetReservationByID(ctx context.Context, id uint) (*servicedto.Reservation, error) {
	r, ok := f.reservations[id]
	if !ok {
		return nil, nil
	}
	copy := r
	return &copy, nil
}

func (f *controllerFakeReservationClient) UpdateReservationStatus(ctx context.Context, id uint, status string) (*servicedto.Reservation, error) {
	r, ok := f.reservations[id]
	if !ok {
		return nil, nil
	}
	r.Status = status
	r.UpdatedAt = time.Now()
	f.reservations[id] = r
	copy := r
	return &copy, nil
}

func (f *controllerFakeReservationClient) ListReservationsByDate(ctx context.Context, date time.Time, status *string) ([]servicedto.Reservation, error) {
	var list []servicedto.Reservation
	for _, r := range f.reservations {
		if !r.Date.Equal(date) {
			continue
		}
		if status != nil && r.Status != *status {
			continue
		}
		list = append(list, r)
	}
	return list, nil
}
