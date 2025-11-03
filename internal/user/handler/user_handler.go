package userhandler

import (
	"errors"

	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/internal/user"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
	"github.com/codepnw/mini-ecommerce/pkg/response"
	"github.com/codepnw/mini-ecommerce/pkg/validate"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	uc userusecase.UserUsecase
}

func NewUserHandler(uc userusecase.UserUsecase) *userHandler {
	return &userHandler{uc: uc}
}

func (h *userHandler) Register(c *gin.Context) {
	req := new(UserCreateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input := &user.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	result, err := h.uc.Register(c, input)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExists) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err)
		return
	}
	response.Created(c, result)
}

func (h *userHandler) Login(c *gin.Context) {
	req := new(UserLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input := &user.User{
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := h.uc.Login(c, input)
	if err != nil {
		if errors.Is(err, errs.ErrUserCredentials) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, "login successfully", result)
}

func (h *userHandler) RefreshToken(c *gin.Context) {
	req := new(RefreshTokenReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	newToken, err := h.uc.RefreshToken(c, req.RefreshToken)
	if err != nil {
		switch err {
		case errs.ErrTokenRevoked:
			response.Unauthorized(c, err.Error())
			return
		case errs.ErrTokenExpires:
			response.Unauthorized(c, err.Error())
			return
		case errs.ErrTokenNotFound:
			response.Unauthorized(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.OK(c, "generate new token successfully", newToken)
}

func (h *userHandler) Logout(c *gin.Context) {
	req := new(RefreshTokenReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.uc.Logout(c, req.RefreshToken); err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.NoContent(c)
}
