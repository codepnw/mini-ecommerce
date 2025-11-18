package orderusecase

type OrderView struct {
	ID        int64            `json:"id"`
	Status    string           `json:"status"`
	Total     float64          `json:"total"`
	CreatedAt string           `json:"created_at"`
	Items     []*OrderItemView `json:"items"`
}

type OrderItemView struct {
	OrderItemID     int64   `json:"item_id"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `json:"price"`
	Total           float64 `json:"total"`
	ProductID       int64   `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductSKU      string  `json:"product_sku"`
}

type OrderListView struct {
	ID        int64   `json:"id"`
	Status    string  `json:"status"`
	Total     float64 `json:"total"`
	CreatedAt string  `json:"created_at"`
}
