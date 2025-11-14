package middleware

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/pkg/auth"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/codepnw/mini-ecommerce/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	token *jwt.JWTToken
}

func InitAuthMiddleware(token *jwt.JWTToken) (*AuthMiddleware, error) {
	if token == nil {
		return nil, errors.New("jwt token is nil")
	}
	return &AuthMiddleware{token: token}, nil
}

func (a *AuthMiddleware) AuthorizedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "auth header is missing")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid token format")
			return
		}

		claims, err := a.token.VerifyAccessToken(parts[1])
		if err != nil {
			log.Printf("verify token failed: %v", err)
			response.Unauthorized(c, "verify token failed")
			return
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, consts.UserClaimsKey, claims)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (a *AuthMiddleware) SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID int64 = 0

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := a.token.VerifyAccessToken(parts[1])
				if err == nil {
					userID = claims.ID
				}
				// Skip all error
			}
		}

		// SessionID
		sessionID := auth.GetOrCreateSessionID(c)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, consts.UserIDKey, userID)
		ctx = context.WithValue(ctx, consts.SessionIDKey, sessionID)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
