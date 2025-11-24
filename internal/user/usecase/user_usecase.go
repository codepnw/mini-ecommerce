package userusecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/codepnw/mini-ecommerce/internal/user"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/pkg/database"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/codepnw/mini-ecommerce/pkg/password"
	"github.com/codepnw/mini-ecommerce/pkg/validate"
)

type UserUsecase interface {
	Register(ctx context.Context, input *user.User) (*TokenResponse, error)
	Login(ctx context.Context, input *user.User) (*TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Logout(ctx context.Context, token string) error

	GetUser(ctx context.Context, userID int64) (*user.User, error)
}

type UserUsecaseConfig struct {
	Repo  userrepository.UserRepository `validate:"required"`
	Token *jwt.JWTToken                 `validate:"required"`
	Tx    database.TxManager            `validate:"required"`
}

type userUsecase struct {
	repo  userrepository.UserRepository
	token *jwt.JWTToken
	tx    database.TxManager
}

func NewUserUsecase(cfg *UserUsecaseConfig) (UserUsecase, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}
	return &userUsecase{
		repo:  cfg.Repo,
		token: cfg.Token,
		tx:    cfg.Tx,
	}, nil
}

func (u *userUsecase) Register(ctx context.Context, input *user.User) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Hashed Password
	hashedPassword, err := password.HashedPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.Password = hashedPassword

	var response *TokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Insert User
		userData, err := u.repo.Insert(ctx, tx, input)
		if err != nil {
			return err
		}

		// Generate Token
		resp, err := u.tokenGenerate(userData)
		if err != nil {
			return err
		}

		// Save Token
		inputAuth := u.inputAuth(userData.ID, resp.RefreshToken)
		if err := u.repo.SaveRefreshToken(ctx, tx, inputAuth); err != nil {
			return err
		}

		response = resp
		return nil
	})
	return response, err
}

func (u *userUsecase) Login(ctx context.Context, input *user.User) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Find User
	userData, err := u.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrUserCredentials
		}
		return nil, err
	}

	// Compare Password
	err = password.ComparePassword(userData.Password, input.Password)
	if err != nil {
		return nil, errs.ErrUserCredentials
	}

	var response *TokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Generate Token
		resp, err := u.tokenGenerate(userData)
		if err != nil {
			return err
		}

		// Save Token
		inputAuth := u.inputAuth(userData.ID, resp.RefreshToken)
		if err := u.repo.SaveRefreshToken(ctx, tx, inputAuth); err != nil {
			return err
		}

		response = resp
		return nil
	})
	return response, err
}

func (u *userUsecase) RefreshToken(ctx context.Context, token string) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Validate Token
	userID, err := u.repo.ValidateRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Find User
	userData, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response *TokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Revoked Token
		if err := u.repo.RevokedRefreshToken(ctx, tx, token); err != nil {
			return err
		}

		// Generate Token
		resp, err := u.tokenGenerate(userData)
		if err != nil {
			return err
		}

		// Save Token
		inputAuth := u.inputAuth(userData.ID, resp.RefreshToken)
		if err := u.repo.SaveRefreshToken(ctx, tx, inputAuth); err != nil {
			return err
		}

		response = resp
		return nil
	})
	return response, err
}

func (u *userUsecase) Logout(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	err := u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := u.repo.RevokedRefreshToken(ctx, tx, refreshToken); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (u *userUsecase) GetUser(ctx context.Context, userID int64) (*user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userData, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *userUsecase) tokenGenerate(input *user.User) (*TokenResponse, error) {
	accessToken, err := u.token.GenerateAccessToken(input)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.token.GenerateRefreshToken(input)
	if err != nil {
		return nil, err
	}

	response := &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}

func (u *userUsecase) inputAuth(userID int64, refreshToken string) *user.Auth {
	return &user.Auth{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(consts.RefreshTokenDuration),
	}
}
