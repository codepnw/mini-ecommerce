package carthandler

type AddItemInput struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"gt=0"`
}
