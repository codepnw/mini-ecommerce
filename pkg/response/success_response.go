package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{
		"response": SuccessResponse{
			Code:    http.StatusCreated,
			Message: "created successfully",
			Data:    data,
		},
	})
}

func OK(c *gin.Context, message string, data any) {
	if message == "" {
		message = "successfully"
	}
	c.JSON(http.StatusOK, gin.H{
		"response": SuccessResponse{
			Code:    http.StatusOK,
			Message: message,
			Data:    data,
		},
	})
}

func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
