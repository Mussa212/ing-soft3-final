package service

import (
	"context"
	"testing"
	"time"

	servicedto "vesuvio/internal/dto/service"
)

func TestCreateReservationValid(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	out, err := svc.CreateReservation(ctx, servicedto.CreateReservationInput{
		UserID: 1,
		Date:   "2025-12-01",
		Time:   "20:30",
		People: 4,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Reservation.Status != servicedto.StatusPending {
		t.Fatalf("expected status pending, got %s", out.Reservation.Status)
	}
}

func TestCreateReservationInvalidPeople(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	_, err := svc.CreateReservation(ctx, servicedto.CreateReservationInput{
		UserID: 1,
		Date:   "2025-12-01",
		Time:   "20:30",
		People: 0,
	})
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateReservationInvalidDate(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	_, err := svc.CreateReservation(ctx, servicedto.CreateReservationInput{
		UserID: 1,
		Date:   "bad-date",
		Time:   "20:30",
		People: 2,
	})
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestListUserReservationsFilter(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()
	date := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID: 1, Date: date, Time: "20:30", People: 2, Status: servicedto.StatusPending,
	})
	client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID: 1, Date: date, Time: "21:00", People: 2, Status: servicedto.StatusConfirmed,
	})

	status := servicedto.StatusConfirmed
	res, err := svc.ListUserReservations(ctx, servicedto.ListUserReservationsInput{
		UserID: 1,
		Status: &status,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Status != servicedto.StatusConfirmed {
		t.Fatalf("unexpected list result: %+v", res)
	}
}

func TestListUserReservationsInvalidStatus(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	status := "weird"
	_, err := svc.ListUserReservations(ctx, servicedto.ListUserReservationsInput{
		UserID: 1,
		Status: &status,
	})
	if err != ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestCancelReservationForbidden(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	res, _ := client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID:  1,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "20:30",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})

	_, err := svc.CancelReservation(ctx, servicedto.CancelReservationInput{
		UserID:        2,
		ReservationID: res.ID,
	})
	if err != ErrForbiddenReservation {
		t.Fatalf("expected ErrForbiddenReservation, got %v", err)
	}
}

func TestCancelReservationInvalidAndMissing(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	_, err := svc.CancelReservation(ctx, servicedto.CancelReservationInput{})
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	_, err = svc.CancelReservation(ctx, servicedto.CancelReservationInput{
		UserID:        1,
		ReservationID: 123,
	})
	if err != ErrReservationNotFound {
		t.Fatalf("expected ErrReservationNotFound, got %v", err)
	}
}

func TestCancelReservationAlreadyCancelled(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	res, _ := client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID:  1,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "20:30",
		People:  2,
		Status:  servicedto.StatusCancelled,
		Comment: nil,
	})

	out, err := svc.CancelReservation(ctx, servicedto.CancelReservationInput{
		UserID:        1,
		ReservationID: res.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error on already cancelled: %v", err)
	}
	if out.Status != servicedto.StatusCancelled {
		t.Fatalf("expected cancelled status to stay cancelled, got %s", out.Status)
	}
}

func TestListUserReservationsMissingUser(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	_, err := svc.ListUserReservations(ctx, servicedto.ListUserReservationsInput{
		UserID: 0,
	})
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAdminListReservationsValidation(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	_, err := svc.AdminListReservations(ctx, servicedto.AdminListReservationsInput{})
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	status := "weird"
	_, err = svc.AdminListReservations(ctx, servicedto.AdminListReservationsInput{
		Date:   "2025-01-01",
		Status: &status,
	})
	if err != ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}

	_, err = svc.AdminListReservations(ctx, servicedto.AdminListReservationsInput{
		Date: "bad-date",
	})
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput for bad date, got %v", err)
	}
}

func TestAdminCancelReservationUnauthorized(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	res, _ := client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID:  1,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "21:00",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})

	_, err := svc.AdminCancelReservation(ctx, servicedto.User{ID: 2, IsAdmin: false}, res.ID)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCancelReservationSuccess(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	res, _ := client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID:  1,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "20:30",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})

	cancelled, err := svc.CancelReservation(ctx, servicedto.CancelReservationInput{
		UserID:        1,
		ReservationID: res.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cancelled.Status != servicedto.StatusCancelled {
		t.Fatalf("expected cancelled, got %s", cancelled.Status)
	}
}

func TestAdminConfirmAndCancel(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()
	admin := servicedto.User{ID: 99, IsAdmin: true}

	res, _ := client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID:  1,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "21:00",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})

	confirmed, err := svc.ConfirmReservation(ctx, admin, res.ID)
	if err != nil {
		t.Fatalf("confirm error: %v", err)
	}
	if confirmed.Status != servicedto.StatusConfirmed {
		t.Fatalf("expected confirmed, got %s", confirmed.Status)
	}

	cancelled, err := svc.AdminCancelReservation(ctx, admin, res.ID)
	if err != nil {
		t.Fatalf("cancel error: %v", err)
	}
	if cancelled.Status != servicedto.StatusCancelled {
		t.Fatalf("expected cancelled, got %s", cancelled.Status)
	}
}

func TestConfirmReservationUnauthorized(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()

	res, _ := client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID:  1,
		Date:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		Time:    "21:00",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})

	_, err := svc.ConfirmReservation(ctx, servicedto.User{ID: 2, IsAdmin: false}, res.ID)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestConfirmReservationNotFound(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()
	admin := servicedto.User{ID: 99, IsAdmin: true}

	_, err := svc.ConfirmReservation(ctx, admin, 999)
	if err != ErrReservationNotFound {
		t.Fatalf("expected ErrReservationNotFound, got %v", err)
	}
}

func TestAdminListReservationsWithStatus(t *testing.T) {
	client := newFakeReservationClient()
	svc := NewReservationService(client)
	ctx := context.Background()
	date := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID: 1, Date: date, Time: "20:00", People: 2, Status: servicedto.StatusPending,
	})
	client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID: 2, Date: date, Time: "21:00", People: 3, Status: servicedto.StatusConfirmed,
	})

	status := servicedto.StatusConfirmed
	res, err := svc.AdminListReservations(ctx, servicedto.AdminListReservationsInput{
		Date:   "2025-12-01",
		Status: &status,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Status != servicedto.StatusConfirmed {
		t.Fatalf("unexpected admin list result: %+v", res)
	}
}

// fakeReservationClient is an in-memory reservation store for tests.
type fakeReservationClient struct {
	reservations map[uint]servicedto.Reservation
	nextID       uint
}

func newFakeReservationClient() *fakeReservationClient {
	return &fakeReservationClient{
		reservations: make(map[uint]servicedto.Reservation),
		nextID:       1,
	}
}

func (f *fakeReservationClient) CreateReservation(ctx context.Context, params servicedto.CreateReservationParams) (*servicedto.Reservation, error) {
	id := f.nextID
	f.nextID++
	now := time.Now()
	res := servicedto.Reservation{
		ID:        id,
		UserID:    params.UserID,
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

func (f *fakeReservationClient) ListReservationsByUser(ctx context.Context, userID uint, status *string) ([]servicedto.Reservation, error) {
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

func (f *fakeReservationClient) GetReservationByID(ctx context.Context, id uint) (*servicedto.Reservation, error) {
	r, ok := f.reservations[id]
	if !ok {
		return nil, nil
	}
	copy := r
	return &copy, nil
}

func (f *fakeReservationClient) UpdateReservationStatus(ctx context.Context, id uint, status string) (*servicedto.Reservation, error) {
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

func (f *fakeReservationClient) ListReservationsByDate(ctx context.Context, date time.Time, status *string) ([]servicedto.Reservation, error) {
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
