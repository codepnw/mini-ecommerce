package helper

import (
	"strconv"

	"github.com/codepnw/mini-ecommerce/internal/user"
	"github.com/gin-gonic/gin"
)

func GetParamInt(c *gin.Context, key string) (int64, error) {
	result, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func GetCurrentUser(c *gin.Context) (*user.User, error) {
	// TODO: Get From Context later
	return &user.User{
		ID: 10,
	}, nil
}
