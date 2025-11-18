package orderrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/mini-ecommerce/internal/order"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/pkg/database"
)

//go:generate mockgen -source=order_repository.go -destination=mock_order_repository.go -package=orderrepository

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, input *order.Order) (int64, error)
	CreateOrderItem(ctx context.Context, tx *sql.Tx, input *order.OrderItem) error
	GetOrder(ctx context.Context, orderID int64) (*order.Order, error)
	GetOrderItems(ctx context.Context, exec database.DBExec, orderID int64) ([]*OrderItemDetail, error)
	GetMyOrders(ctx context.Context, userID int64) ([]*order.Order, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, orderID int64, status string) error
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

func (r *orderRepository) GetOrder(ctx context.Context, orderID int64) (*order.Order, error) {
	query := `
		SELECT id, user_id, total, status, created_at, updated_at
		FROM orders WHERE id = $1 LIMIT 1
	`
	o := new(order.Order)
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&o.ID,
		&o.UserID,
		&o.Total,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return o, nil
}

type OrderItemDetail struct {
	ID              int64   `json:"id"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `json:"price"`
	ProductID       int64   `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductSKU      string  `json:"product_sku"`
}

func (r *orderRepository) GetOrderItems(ctx context.Context, exec database.DBExec, orderID int64) ([]*OrderItemDetail, error) {
	query := `
		SELECT oi.id, oi.product_id, oi.quantity, oi.price, p.name, p.sku
		FROM order_items oi
		INNER JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
	`
	rows, err := exec.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*OrderItemDetail, 0)
	for rows.Next() {
		i := new(OrderItemDetail)
		err = rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.Quantity,
			&i.PriceAtPurchase,
			&i.ProductName,
			&i.ProductSKU,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *orderRepository) GetMyOrders(ctx context.Context, userID int64) ([]*order.Order, error) {
	query := `
		SELECT id, user_id, total, status, created_at
		FROM orders WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*order.Order, 0)
	for rows.Next() {
		o := new(order.Order)
		err = rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Total,
			&o.Status,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, orderID int64, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, status, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrOrderNotFound
		}
	}
	return nil
}
