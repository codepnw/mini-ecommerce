package productrepository

import (
	"database/sql"
	"time"

	"github.com/codepnw/mini-ecommerce/internal/product"
)

type productModel struct {
	ID        int64          `db:"id"`
	Name      string         `db:"name"`
	Price     float64        `db:"price"`
	Stock     int            `db:"stock"`
	SKU       sql.NullString `db:"sku"`
	OwnerID   sql.NullInt64  `db:"owner_id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (r *productRepository) inputToModel(p *product.Product) *productModel {
	nullSKU := sql.NullString{String: p.SKU, Valid: p.SKU != ""}
	nullOwnerID := sql.NullInt64{Int64: p.OwnerID, Valid: p.OwnerID > 0}

	return &productModel{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		SKU:       nullSKU,
		OwnerID:   nullOwnerID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (r *productRepository) modelToDomain(p *productModel) *product.Product {
	return &product.Product{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		SKU:       p.SKU.String,
		OwnerID:   p.OwnerID.Int64,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
