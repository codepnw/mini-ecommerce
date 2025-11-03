package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParamInt(c *gin.Context, key string) (int64, error) {
	result, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
