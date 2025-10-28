package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": ErrorResponse{
			Code:    http.StatusBadRequest,
			Type:    "BAD_REQUEST",
			Message: message,
		},
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": ErrorResponse{
			Code:    http.StatusNotFound,
			Type:    "NOT_FOUND",
			Message: message,
		},
	})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": ErrorResponse{
			Code:    http.StatusInternalServerError,
			Type:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		},
	})
}
