package producthandler

type ProductCreateReq struct {
	Name  string  `json:"name" binding:"required,min=2"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int     `json:"stock" binding:"gt=0"`
	SKU   string  `json:"sku" binding:"required,min=2,max=20"`
}

type ProductUpdateReq struct {
	Name  *string  `json:"name,omitempty" binding:"omitempty,min=2"`
	Price *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
	Stock *int     `json:"stock,omitempty" binding:"omitempty,gt=0"`
	SKU   *string  `json:"sku,omitempty" binding:"omitempty,min=2,max=20"`
}
