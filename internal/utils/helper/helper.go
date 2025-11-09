package helper

import (
	"context"
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

func GetUserIDFromCtx(ctx context.Context) int64 {
	// TODO: change later
	return 10
}
func GetSessionIDFromCtx(ctx context.Context) string {
	// TODO: change later
	return "mock_session_id"
}
