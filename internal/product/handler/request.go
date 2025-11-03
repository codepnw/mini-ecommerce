package producthandler

type ProductCreateReq struct {
	Name  string  `json:"name" validate:"required,min=2"`
	Price float64 `json:"price" validate:"required,gt=0"`
	Stock int     `json:"stock" validate:"gt=0"`
}

type ProductUpdateReq struct {
	Name  *string  `json:"name,omitempty" validate:"omitempty,min=2"`
	Price *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock *int    `json:"stock,omitempty" validate:"omitempty,gt=0"`
}
