package service

import "errors"

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserNotFound         = errors.New("user not found")
	ErrReservationNotFound  = errors.New("reservation not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrInvalidInput         = errors.New("invalid input")
	ErrForbiddenReservation = errors.New("user cannot modify this reservation")
)
