package carthandler

type AddItemReq struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"gt=0"`
}

type UpdateItemQuantityReq struct {
	NewQuantity int   `json:"new_quantity"`
}
