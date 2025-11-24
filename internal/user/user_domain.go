package user

import "time"

type RoleType string

const (
	RoleUser   RoleType = "user"
	RoleSeller RoleType = "seller"
	RoleAdmin  RoleType = "admin"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Auth struct {
	ID           int64
	UserID       int64
	RefreshToken string
	Revoked      bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
