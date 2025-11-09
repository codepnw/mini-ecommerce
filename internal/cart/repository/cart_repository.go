package cartrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/mini-ecommerce/internal/cart"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
)

type CartRepository interface {
	GetOrCreateActiveCart(ctx context.Context, userID sql.NullInt64, sessionID string) (*cart.Cart, error)
	UpsertItem(ctx context.Context, tx *sql.Tx, item *cart.CartItem) error
	GetCartItems(ctx context.Context, cartID string) ([]*CartItemDB, error)
	UpdateItemQuantity(ctx context.Context, tx *sql.Tx, cartID string, cartItemID int64, quantity int) error
	RemoveItem(ctx context.Context, tx *sql.Tx, cartID string, cartItemID int64) error
	ClearCart(ctx context.Context, cartID string) error
	GetCartItemDetails(ctx context.Context, cartItemID int64, cartID string) (*cart.CartItem, error)
	GetCartItemForUpdate(ctx context.Context, tx *sql.Tx, cartItemID int64, cartID string) (*cart.CartItem, error)
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetOrCreateActiveCart(ctx context.Context, userID sql.NullInt64, sessionID string) (*cart.Cart, error) {
	var query string
	var args []any
	c := new(cart.Cart)

	if userID.Valid {
		query = `SELECT cart_id, user_id, session_id, status FROM carts WHERE user_id = $1 AND status = 'active' LIMIT 1`
		args = append(args, userID.Int64)
	} else {
		query = `SELECT cart_id, user_id, session_id, status FROM carts WHERE session_id = $1 AND status = 'guest' LIMIT 1`
		args = append(args, sessionID)
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&c.CartID,
		&c.UserID,
		&c.SessionID,
		&c.Status,
	)
	if err == nil {
		// Found Cart
		return c, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Other Error
		return nil, err
	}

	// Not Found (err == sql.ErrNoRows) -> Create Cart
	var insertQuery string
	var insertArgs []any

	if userID.Valid {
		insertQuery = `INSERT INTO carts (user_id, status) VALUES ($1, 'active') RETURNING cart_id, user_id, session_id, status`
		insertArgs = append(insertArgs, userID.Int64)
	} else {
		insertQuery = `INSERT INTO carts (session_id, status) VALUES ($1, 'guest') RETURNING cart_id, user_id, session_id, status`
		insertArgs = append(insertArgs, sessionID)
	}

	newCart := new(cart.Cart)
	err = r.db.QueryRowContext(ctx, insertQuery, insertArgs...).Scan(
		&newCart.CartID,
		&newCart.UserID,
		&newCart.SessionID,
		&newCart.Status,
	)
	if err != nil {
		return nil, err
	}
	return newCart, nil
}

func (r *cartRepository) UpsertItem(ctx context.Context, tx *sql.Tx, item *cart.CartItem) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity, price_at_add)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET
			quantity = cart_items.quantity + EXCLUDED.quantity,
			updated_at = NOW()
	`
	_, err := tx.ExecContext(
		ctx,
		query,
		item.CartID,
		item.ProductID,
		item.Quantity,
		item.PriceAtAdd,
	)
	if err != nil {
		return err
	}
	return nil
}

type CartItemDB struct {
	CartItemID string
	ProductID  int64
	Quantity   int
	PriceAtAdd float64
	// From Products
	Name  string
	Price float64
	Stock int
	SKU   string
}

func (r *cartRepository) GetCartItems(ctx context.Context, cartID string) ([]*CartItemDB, error) {
	query := `
		SELECT ci.id, ci.product_id, ci.quantity, ci.price_at_add, p.name, p.price, p.stock, p.sku
		FROM cart_items ci
		INNER JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
		ORDER BY ci.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*CartItemDB, 0)

	for rows.Next() {
		item := new(CartItemDB)
		err = rows.Scan(
			&item.CartItemID,
			&item.ProductID,
			&item.Quantity,
			&item.PriceAtAdd,
			&item.Name,
			&item.Price,
			&item.Stock,
			&item.SKU,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartRepository) UpdateItemQuantity(ctx context.Context, tx *sql.Tx, cartID string, cartItemID int64, quantity int) error {
	query := `
		UPDATE cart_items SET quantity = $1, updated_at = NOW()
		WHERE id = $2 AND cart_id = $3
	`
	_, err := tx.ExecContext(ctx, query, quantity, cartItemID, cartID)
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) RemoveItem(ctx context.Context, tx *sql.Tx, cartID string, cartItemID int64) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1 AND id = $2`
	_, err := tx.ExecContext(ctx, query, cartID, cartItemID)
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) ClearCart(ctx context.Context, cartID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) GetCartItemDetails(ctx context.Context, cartItemID int64, cartID string) (*cart.CartItem, error) {
	query := `
		SELECT id, product_id, quantity, cart_id
		FROM cart_items
		WHERE id = $1 AND cart_id = $2 LIMIT 1
	`
	item := new(cart.CartItem)
	err := r.db.QueryRowContext(ctx, query, cartItemID, cartID).Scan(
		&item.ID,
		&item.ProductID,
		&item.Quantity,
		&item.CartID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrItemNotInCart
		}
		return nil, err
	}
	return item, nil
}

func (r *cartRepository) GetCartItemForUpdate(ctx context.Context, tx *sql.Tx, cartItemID int64, cartID string) (*cart.CartItem, error) {
	query := `
		SELECT id, product_id, quantity, cart_id
		FROM cart_items
		WHERE id = $1 AND cart_id = $2 LIMIT 1
		FOR UPDATE
	`
	item := new(cart.CartItem)
	err := tx.QueryRowContext(ctx, query, cartItemID, cartID).Scan(
		&item.ID,
		&item.ProductID,
		&item.Quantity,
		&item.CartID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrItemNotInCart
		}
		return nil, err
	}
	return item, nil
}
