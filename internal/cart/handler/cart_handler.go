package carthandler

import (
	cartusecase "github.com/codepnw/mini-ecommerce/internal/cart/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/internal/utils/helper"
	"github.com/codepnw/mini-ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
)

type cartHandler struct {
	uc cartusecase.CartUsecase
}

func NewCartHandler(uc cartusecase.CartUsecase) *cartHandler {
	return &cartHandler{uc: uc}
}

func (h *cartHandler) AddItemToCart(c *gin.Context) {
	req := new(AddItemReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.AddItemToCart(c.Request.Context(), req.ProductID, req.Quantity)
	if err != nil {
		switch err {
		case errs.ErrProductNotEnough:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrProductNotFound:
			response.NotFound(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.OK(c, "", result)
}

func (h *cartHandler) GetCart(c *gin.Context) {
	result, err := h.uc.GetCart(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, "", result)
}

func (h *cartHandler) UpdateItemQuantity(c *gin.Context) {
	id, err := helper.GetParamInt(c, consts.CartItemID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	req := new(UpdateItemQuantityReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.UpdateItemQuantity(c.Request.Context(), id, req.NewQuantity)
	if err != nil {
		switch err {
		case errs.ErrInvalidQuantity:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrItemNotInCart:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrProductNotEnough:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrProductNotFound:
			response.NotFound(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.OK(c, "new quantity updated", result)
}

func (h *cartHandler) RemoveItemFromCart(c *gin.Context) {
	id, err := helper.GetParamInt(c, consts.CartItemID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.uc.RemoveItemFromCart(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, "item remove", result)
}
