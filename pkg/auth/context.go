package auth

import (
	"context"
	"errors"
	"time"

	"github.com/codepnw/mini-ecommerce/internal/user"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetOrCreateSessionID(c *gin.Context) string {
	const sessionCookieName = "session_id"

	sessionID, err := c.Cookie(sessionCookieName)
	if err == nil {
		return sessionID
	}

	newSessionID := uuid.NewString()
	maxAge := int(time.Second * 24 * 365) // 1 year

	c.SetCookie(
		sessionCookieName,
		newSessionID,
		maxAge,
		"/",
		"",
		true,
		true,
	)
	return newSessionID
}

func GetCurrentUser(ctx context.Context) (*user.User, error) {
	claims, ok := ctx.Value(consts.UserClaimsKey).(*jwt.UserClaims)
	if !ok {
		return nil, errors.New("no user")
	}

	usr := &user.User{
		ID:    claims.ID,
		Email: claims.Email,
		Role:  claims.Role,
	}
	return usr, nil
}

func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(consts.UserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

func GetSessionID(ctx context.Context) string {
	sessionID, ok := ctx.Value(consts.SessionIDKey).(string)
	if !ok {
		return ""
	}
	return sessionID
}

func SetCurrentUser(ctx context.Context, claims *jwt.UserClaims) context.Context {
	return context.WithValue(ctx, consts.UserClaimsKey, claims)
}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, consts.UserIDKey, userID)
}
