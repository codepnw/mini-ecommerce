package productrepository

import (
	"time"

	"github.com/codepnw/mini-ecommerce/internal/product"
)

type productModel struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Price     float64   `db:"price"`
	Stock     int       `db:"stock"`
	SKU       string    `db:"sku"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *productRepository) inputToModel(p *product.Product) *productModel {
	return &productModel{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		SKU:       p.SKU,
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
		SKU:       p.SKU,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
