package orderrepository

import (
	"context"
	"database/sql"

	"github.com/codepnw/mini-ecommerce/internal/order"
)

//go:generate mockgen -source=order_repository.go -destination=mock_order_repository.go -package=orderrepository

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, input *order.Order) (int64, error)
	CreateOrderItem(ctx context.Context, tx *sql.Tx, input *order.OrderItem) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, tx *sql.Tx, input *order.Order) (int64, error) {
	query := `
		INSERT INTO orders (user_id, total, status)
		VALUES ($1, $2, $3) RETURNING id
	`
	var orderID int64
	err := tx.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.Total,
		input.Status,
	).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *orderRepository) CreateOrderItem(ctx context.Context, tx *sql.Tx, input *order.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, price, quantity)
		VALUES ($1, $2, $3, $4)
	`
	_, err := tx.ExecContext(
		ctx,
		query,
		input.OrderID,
		input.ProductID,
		input.PriceAtPurchase,
		input.Quantity,
 	)
	return err
}
