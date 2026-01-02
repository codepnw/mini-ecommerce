package orderhandler

type UpdateStatusReq struct {
	Status string `json:"status" binding:"required,oneof=paid shipped cancelled completed"`
}
