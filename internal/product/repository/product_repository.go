package productrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codepnw/mini-ecommerce/internal/product"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
)

//go:generate mockgen -source=product_repository.go -destination=mock_product_repository.go -package=productrepository

type ProductRepository interface {
	Insert(ctx context.Context, input *product.Product) (*product.Product, error)
	FindByID(ctx context.Context, id int64) (*product.Product, error)
	FindByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*product.Product, error)
	List(ctx context.Context) ([]*product.Product, error)
	Update(ctx context.Context, input *product.Product) (*product.Product, error)
	Delete(ctx context.Context, id int64) error

	SKUExists(ctx context.Context, sku string) (bool, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Insert(ctx context.Context, input *product.Product) (*product.Product, error) {
	m := r.inputToModel(input)
	query := `
		INSERT INTO products (name, price, stock, sku, owner_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		m.Name,
		m.Price,
		m.Stock,
		m.SKU,
		m.OwnerID,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *productRepository) FindByID(ctx context.Context, id int64) (*product.Product, error) {
	pd := new(product.Product)
	p := r.inputToModel(pd)

	query := `
		SELECT id, name, price, stock, sku, owner_id, created_at, updated_at
		FROM products WHERE id = $1 LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.SKU,
		&p.OwnerID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}
	return r.modelToDomain(p), nil
}

func (r *productRepository) FindByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*product.Product, error) {
	pd := new(product.Product)
	p := r.inputToModel(pd)

	query := `
		SELECT id, name, price, stock, sku, owner_id, created_at, updated_at
		FROM products WHERE id = $1 LIMIT 1 FOR UPDATE
	`
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.SKU,
		&p.OwnerID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}
	return r.modelToDomain(p), nil
}

func (r *productRepository) List(ctx context.Context) ([]*product.Product, error) {
	query := `
		SELECT id, name, price, stock, sku, owner_id, created_at, updated_at
		FROM products
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var products []*product.Product
	for rows.Next() {
		p := new(product.Product)
		if err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.SKU,
			&p.OwnerID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *productRepository) Update(ctx context.Context, input *product.Product) (*product.Product, error) {
	var (
		sb      strings.Builder
		columns []string
		values  []any
		idx     = 1
	)

	sb.WriteString("UPDATE products SET updated_at = NOW()")
	if input.Name != "" {
		columns = append(columns, fmt.Sprintf("name = $%d", idx))
		values = append(values, input.Name)
		idx++
	}
	if input.Price != 0 {
		columns = append(columns, fmt.Sprintf("price = $%d", idx))
		values = append(values, input.Price)
		idx++
	}
	if input.Stock != 0 {
		columns = append(columns, fmt.Sprintf("stock = $%d", idx))
		values = append(values, input.Stock)
		idx++
	}
	if input.SKU != "" {
		columns = append(columns, fmt.Sprintf("sku = $%d", idx))
		values = append(values, input.SKU)
		idx++
	}

	if len(columns) > 0 {
		sb.WriteString(", ") // for updated_at
		sb.WriteString(strings.Join(columns, ", "))
	}

	sb.WriteString(fmt.Sprintf(" WHERE id = $%d", idx))
	values = append(values, input.ID)
	idx++

	sb.WriteString(" RETURNING id, name, price, stock, sku, owner_id, created_at, updated_at")

	query := sb.String()
	log.Println(query)

	p := new(product.Product)
	err := r.db.QueryRowContext(ctx, query, values...).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.SKU,
		&p.OwnerID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductNotFound
	}
	return nil
}

func (r *productRepository) SKUExists(ctx context.Context, sku string) (bool, error) {
	var exists int
	query := `SELECT 1 FROM products WHERE sku = $1 LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, sku).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
