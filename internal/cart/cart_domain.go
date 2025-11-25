package cart

import (
	"database/sql"
	"time"
)

type status string

const (
	StatusActive  status = "active"
	StatusGuest   status = "guest"
	StatusSaved   status = "saved"
	StatusOrdered status = "ordered"
)

type Cart struct {
	ID        string         `json:"id"`
	UserID    sql.NullInt64  `json:"user_id"`
	SessionID sql.NullString `json:"sesstion_id"`
	Status    status         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type CartItem struct {
	ID         int64     `json:"id"`
	CartID     string    `json:"cart_id"`
	ProductID  int64     `json:"product_id"`
	Quantity   int       `json:"quantity"`
	PriceAtAdd float64   `json:"price_at_add"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
