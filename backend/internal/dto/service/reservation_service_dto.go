package servicedto

import "time"

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
)

// Reservation is the service-level representation.
type Reservation struct {
	ID        uint
	UserID    uint
	User      *User
	Date      time.Time
	Time      string
	People    int
	Comment   *string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateReservationInput carries data for creating a reservation.
type CreateReservationInput struct {
	UserID  uint
	Date    string
	Time    string
	People  int
	Comment *string
}

type CreateReservationOutput struct {
	Reservation Reservation
}

type ListUserReservationsInput struct {
	UserID uint
	Status *string
}

type CancelReservationInput struct {
	UserID        uint
	ReservationID uint
}

type AdminListReservationsInput struct {
	Date   string
	Status *string
}

// CreateReservationParams used by the client layer when persisting.
type CreateReservationParams struct {
	UserID  uint
	Date    time.Time
	Time    string
	People  int
	Comment *string
	Status  string
}
