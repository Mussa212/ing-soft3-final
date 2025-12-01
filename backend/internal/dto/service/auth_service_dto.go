package servicedto

import "time"

// User is the service-level representation.
type User struct {
	ID        uint
	Name      string
	Email     string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserWithPassword includes sensitive info for internal use only.
type UserWithPassword struct {
	User
	PasswordHash string
}

// RegisterUserInput carries register data.
type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserOutput struct {
	User User
}

// LoginUserInput carries login data.
type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	User User
}

// CreateUserParams is used by the client to persist a user.
type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
	IsAdmin      bool
}
