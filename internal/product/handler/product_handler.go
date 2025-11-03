package producthandler

import (
	"errors"

	"github.com/codepnw/mini-ecommerce/internal/product"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/internal/utils/helper"
	"github.com/codepnw/mini-ecommerce/pkg/response"
	"github.com/codepnw/mini-ecommerce/pkg/validate"
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
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input := &product.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	resp, err := h.uc.Create(c, input)
	if err != nil {
		switch err {
		case errs.ErrProductPriceInvalid:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrProductStockInvalid:
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

	resp, err := h.uc.GetByID(c, id)
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

func (h *productHandler) GetAll(c *gin.Context) {
	// TODO: filter

	resp, err := h.uc.GetAll(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, "", resp)
}

func (h *productHandler) Update(c *gin.Context) {
	id, err := h.getParamID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	req := new(ProductUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var hasUpdate bool

	input := &product.Product{
		ID: id,
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

	if !hasUpdate {
		response.BadRequest(c, errs.ErrNoFieldsToUpdate.Error())
		return
	}

	resp, err := h.uc.Update(c, input)
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

func (h *productHandler) Delete(c *gin.Context) {
	id, err := h.getParamID(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err = h.uc.Delete(c, id); err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err)
		return
	}
	response.NoContent(c)
}

func (h *productHandler) getParamID(c *gin.Context) (int64, error) {
	return helper.GetParamInt(c, consts.ParamProductID)
}
