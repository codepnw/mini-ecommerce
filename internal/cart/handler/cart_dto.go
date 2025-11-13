package carthandler

type AddItemReq struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"gt=0"`
}

type UpdateItemQuantityReq struct {
	NewQuantity int   `json:"new_quantity"`
}
