package controllerdto

// AdminReservationResponse includes reservation plus user info.
type AdminReservationResponse struct {
	ID        uint          `json:"id"`
	User      AdminUserInfo `json:"user"`
	Date      string        `json:"date"`
	Time      string        `json:"time"`
	People    int           `json:"people"`
	Comment   *string       `json:"comment,omitempty"`
	Status    string        `json:"status"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

// AdminUserInfo exposes limited user data in admin responses.
type AdminUserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
