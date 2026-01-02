package producthandler

import (
	"errors"

	"github.com/codepnw/mini-ecommerce/internal/product"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/internal/utils/helper"
	"github.com/codepnw/mini-ecommerce/pkg/auth"
	"github.com/codepnw/mini-ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	uc productusecase.ProductUsecase
}

func NewProductHandler(uc productusecase.ProductUsecase) *productHandler {
	return &productHandler{uc: uc}
}

func (h *productHandler) Create(c *gin.Context) {
	req := new(ProductCreateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get User Context
	userCtx, err := auth.GetCurrentUser(c.Request.Context())
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	input := &product.Product{
		Name:    req.Name,
		Price:   req.Price,
		Stock:   req.Stock,
		SKU:     req.SKU,
		OwnerID: userCtx.ID,
	}
	resp, err := h.uc.Create(c.Request.Context(), input)
	if err != nil {
		switch err {
		case errs.ErrProductPriceInvalid:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrProductStockInvalid:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrProductSKUExists:
			response.BadRequest(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.Created(c, resp)
}

func (h *productHandler) GetByID(c *gin.Context) {
	id, err := h.getParamID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, "", resp)
}

func (h *productHandler) List(c *gin.Context) {
	filter := new(product.ProductFilter)
	if err := c.ShouldBindQuery(filter); err != nil {
		response.BadRequest(c, "invalid filter params")
		return
	}

	resp, err := h.uc.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, "", resp)
}

func (h *productHandler) Update(c *gin.Context) {
	productID, err := h.getParamID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	req := new(ProductUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var hasUpdate bool

	input := &product.Product{
		ID: productID,
	}
	if req.Name != nil {
		input.Name = *req.Name
		hasUpdate = true
	}
	if req.Price != nil {
		input.Price = *req.Price
		hasUpdate = true
	}
	if req.Stock != nil {
		input.Stock = *req.Stock
		hasUpdate = true
	}
	if req.SKU != nil {
		input.SKU = *req.SKU
		hasUpdate = true
	}

	if !hasUpdate {
		response.BadRequest(c, errs.ErrNoFieldsToUpdate.Error())
		return
	}

	resp, err := h.uc.Update(c.Request.Context(), input)
	if err != nil {
		switch err {
		case errs.ErrProductNotFound:
			response.NotFound(c, err.Error())
			return
		case errs.ErrNoPermissions:
			response.Forbidden(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.OK(c, "", resp)
}

func (h *productHandler) Delete(c *gin.Context) {
	productID, err := h.getParamID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err = h.uc.Delete(c.Request.Context(), productID); err != nil {
		switch err {
		case errs.ErrProductNotFound:
			response.NotFound(c, err.Error())
			return
		case errs.ErrNoPermissions:
			response.Forbidden(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.NoContent(c)
}

func (h *productHandler) getParamID(c *gin.Context) (int64, error) {
	return helper.GetParamInt(c, consts.ParamProductID)
}
