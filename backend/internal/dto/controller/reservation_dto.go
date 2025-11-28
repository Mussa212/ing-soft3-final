package controllerdto

// CreateReservationRequest payload for creating a reservation.
type CreateReservationRequest struct {
	Date    string  `json:"date" binding:"required"` // YYYY-MM-DD
	Time    string  `json:"time" binding:"required"` // HH:MM
	People  int     `json:"people" binding:"required"`
	Comment *string `json:"comment,omitempty"`
}

// ReservationResponse basic reservation data for clients.
type ReservationResponse struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	Date      string  `json:"date"`
	Time      string  `json:"time"`
	People    int     `json:"people"`
	Comment   *string `json:"comment,omitempty"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
