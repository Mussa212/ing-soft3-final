package service

import (
	"context"
	"time"

	servicedto "vesuvio/internal/dto/service"
)

// ReservationClient abstracts reservation persistence.
type ReservationClient interface {
	CreateReservation(ctx context.Context, params servicedto.CreateReservationParams) (*servicedto.Reservation, error)
	ListReservationsByUser(ctx context.Context, userID uint, status *string) ([]servicedto.Reservation, error)
	GetReservationByID(ctx context.Context, id uint) (*servicedto.Reservation, error)
	UpdateReservationStatus(ctx context.Context, id uint, status string) (*servicedto.Reservation, error)
	ListReservationsByDate(ctx context.Context, date time.Time, status *string) ([]servicedto.Reservation, error)
}

type ReservationService struct {
	reservationClient ReservationClient
}

func NewReservationService(resClient ReservationClient) *ReservationService {
	return &ReservationService{reservationClient: resClient}
}

func (s *ReservationService) CreateReservation(ctx context.Context, input servicedto.CreateReservationInput) (*servicedto.CreateReservationOutput, error) {
	return nil, nil
}

func (s *ReservationService) ListUserReservations(ctx context.Context, input servicedto.ListUserReservationsInput) ([]servicedto.Reservation, error) {
	if input.UserID == 0 {
		return nil, ErrInvalidInput
	}
	if input.Status != nil && !isValidStatus(*input.Status) {
		return nil, ErrInvalidStatus
	}
	return s.reservationClient.ListReservationsByUser(ctx, input.UserID, input.Status)
}

func (s *ReservationService) CancelReservation(ctx context.Context, input servicedto.CancelReservationInput) (*servicedto.Reservation, error) {
	if input.UserID == 0 || input.ReservationID == 0 {
		return nil, ErrInvalidInput
	}

	res, err := s.reservationClient.GetReservationByID(ctx, input.ReservationID)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrReservationNotFound
	}
	if res.UserID != input.UserID {
		return nil, ErrForbiddenReservation
	}

	if res.Status == servicedto.StatusCancelled {
		return res, nil
	}

	updated, err := s.reservationClient.UpdateReservationStatus(ctx, input.ReservationID, servicedto.StatusCancelled)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *ReservationService) AdminListReservations(ctx context.Context, input servicedto.AdminListReservationsInput) ([]servicedto.Reservation, error) {
	if input.Date == "" {
		return nil, ErrInvalidInput
	}
	if input.Status != nil && !isValidStatus(*input.Status) {
		return nil, ErrInvalidStatus
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, ErrInvalidInput
	}

	return s.reservationClient.ListReservationsByDate(ctx, date, input.Status)
}

func (s *ReservationService) ConfirmReservation(ctx context.Context, admin servicedto.User, reservationID uint) (*servicedto.Reservation, error) {
	if !admin.IsAdmin {
		return nil, ErrUnauthorized
	}
	return s.updateReservationStatus(ctx, reservationID, servicedto.StatusConfirmed)
}

func (s *ReservationService) AdminCancelReservation(ctx context.Context, admin servicedto.User, reservationID uint) (*servicedto.Reservation, error) {
	if !admin.IsAdmin {
		return nil, ErrUnauthorized
	}
	return s.updateReservationStatus(ctx, reservationID, servicedto.StatusCancelled)
}

func (s *ReservationService) updateReservationStatus(ctx context.Context, reservationID uint, status string) (*servicedto.Reservation, error) {
	if reservationID == 0 {
		return nil, ErrInvalidInput
	}
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	res, err := s.reservationClient.GetReservationByID(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrReservationNotFound
	}

	return s.reservationClient.UpdateReservationStatus(ctx, reservationID, status)
}

func isValidStatus(status string) bool {
	switch status {
	case servicedto.StatusPending, servicedto.StatusConfirmed, servicedto.StatusCancelled:
		return true
	default:
		return false
	}
}
