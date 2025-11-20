package orderhandler

type UpdateStatusReq struct {
	Status string `json:"status" validate:"required,oneof=paid shipped cancelled completed"`
}
