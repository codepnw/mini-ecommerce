package order

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusCancelled OrderStatus = "cancelled"
	StatusCompleted OrderStatus = "completed"
)

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID              int64   `json:"id"`
	OrderID         int64   `json:"order_id"`
	ProductID       int64   `json:"product_id"`
	PriceAtPurchase float64 `json:"price_at_purchase"`
	Quantity        int     `json:"quantity"`
}
