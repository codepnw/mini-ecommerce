package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/mini-ecommerce/internal/user"
	"github.com/codepnw/mini-ecommerce/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	secretKey  string
	refreshKey string
}

type UserClaims struct {
	ID    int
	Email string
	Role  string
	*jwt.RegisteredClaims
}

func InitJWT(cfg config.JWTConfig) (*JWTToken, error) {
	if cfg.SecretKey == "" || cfg.RefreshKey == "" {
		return nil, errors.New("jwt key is empty string")
	}
	return &JWTToken{
		secretKey:  cfg.SecretKey,
		refreshKey: cfg.RefreshKey,
	}, nil
}

func (t *JWTToken) GenerateAccessToken(u *user.User) (string, error) {
	duration := time.Hour * 24
	return t.generateToken(t.secretKey, u, duration)
}

func (t *JWTToken) GenerateRefreshToken(u *user.User) (string, error) {
	duration := time.Hour * 24 * 7
	return t.generateToken(t.refreshKey, u, duration)
}

func (t *JWTToken) generateToken(key string, u *user.User, duration time.Duration) (string, error) {
	claims := &UserClaims{
		ID:    u.ID,
		Email: u.Email,
		Role:  "user", // TODO: change later
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "mini-ecommerce",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("signed token failed: %w", err)
	}
	return ss, nil
}
