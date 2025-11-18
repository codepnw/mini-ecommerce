package orderhandler

import (
	orderusecase "github.com/codepnw/mini-ecommerce/internal/order/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/internal/utils/helper"
	"github.com/codepnw/mini-ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	uc orderusecase.OrderUsecase
}

func NewOrderHandler(uc orderusecase.OrderUsecase) *orderHandler {
	return &orderHandler{uc: uc}
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	result, err := h.uc.CreateOrder(c.Request.Context())
	if err != nil {
		switch err {
		case errs.ErrUnauthorized:
			response.Unauthorized(c, err.Error())
			return
		case errs.ErrCartIsEmpty:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrProductNotEnough:
			response.BadRequest(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.OK(c, "", result)
}

func (h *orderHandler) GetOrderDetail(c *gin.Context) {
	orderID, err := helper.GetParamInt(c, consts.ParamOrderID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.GetOrderDetail(c.Request.Context(), orderID)
	if err != nil {
		switch err {
		case errs.ErrUnauthorized:
			response.Unauthorized(c, err.Error())
			return
		case errs.ErrNoPermissions:
			response.Forbidden(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.OK(c, "", result)
}

func (h *orderHandler) GetMyOrders(c *gin.Context) {
	result, err := h.uc.GetMyOrders(c.Request.Context())
	if err != nil {
		switch err {
		case errs.ErrUnauthorized:
			response.Unauthorized(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.OK(c, "", result)
}
