package userusecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/codepnw/mini-ecommerce/internal/consts"
	"github.com/codepnw/mini-ecommerce/internal/errs"
	"github.com/codepnw/mini-ecommerce/internal/user"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	"github.com/codepnw/mini-ecommerce/pkg/database"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/codepnw/mini-ecommerce/pkg/password"
	"github.com/codepnw/mini-ecommerce/pkg/validate"
)

type UserUsecase interface {
	Register(ctx context.Context, input *user.User) (*tokenResponse, error)
	Login(ctx context.Context, input *user.User) (*tokenResponse, error)
	RefreshToken(ctx context.Context, token string) (*tokenResponse, error)
	Logout(ctx context.Context, token string) error
}

type UserUsecaseConfig struct {
	Repo  userrepository.UserRepository `validate:"required"`
	Token *jwt.JWTToken                 `validate:"required"`
	Tx    *database.TxManager           `validate:"required"`
	DB    *sql.DB                       `validate:"required"`
}

type userUsecase struct {
	repo  userrepository.UserRepository
	token *jwt.JWTToken
	tx    *database.TxManager
	db    *sql.DB
}

func NewUserUsecase(cfg *UserUsecaseConfig) (UserUsecase, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}
	return &userUsecase{
		repo:  cfg.Repo,
		token: cfg.Token,
		tx:    cfg.Tx,
		db:    cfg.DB,
	}, nil
}

func (u *userUsecase) Register(ctx context.Context, input *user.User) (*tokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Hashed Password
	hashed, err := password.HashedPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.Password = hashed

	var response *tokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Insert User
		userCreated, err := u.repo.Insert(ctx, tx, input)
		if err != nil {
			return err
		}

		// Generate Token
		resp, err := u.tokenGenerate(userCreated)
		if err != nil {
			return err
		}

		// Save Token
		inputAuth := u.inputAuth(userCreated.ID, resp.RefreshToken)
		if err := u.repo.SaveRefreshToken(ctx, tx, inputAuth); err != nil {
			return err
		}

		response = resp
		return nil
	})
	return response, err
}

func (u *userUsecase) Login(ctx context.Context, input *user.User) (*tokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Find User
	userResult, err := u.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserCredentials
		}
		return nil, err
	}

	// Compare Password
	err = password.ComparePassword(userResult.Password, input.Password)
	if err != nil {
		return nil, errs.ErrUserCredentials
	}

	var response *tokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Generate Token
		resp, err := u.tokenGenerate(userResult)
		if err != nil {
			return err
		}

		// Save Token
		inputAuth := u.inputAuth(userResult.ID, response.RefreshToken)
		if err := u.repo.SaveRefreshToken(ctx, tx, inputAuth); err != nil {
			return err
		}

		response = resp
		return nil
	})
	return response, err
}

func (u *userUsecase) RefreshToken(ctx context.Context, token string) (*tokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Validate Token
	userID, err := u.repo.ValidateRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Find User
	userResult, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	var response *tokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Revoked Token
		if err := u.repo.RevokedRefreshToken(ctx, tx, token); err != nil {
			return err
		}

		// Generate Token
		resp, err := u.tokenGenerate(userResult)
		if err != nil {
			return err
		}

		// Save Token
		inputAuth := u.inputAuth(userResult.ID, resp.RefreshToken)
		if err := u.repo.SaveRefreshToken(ctx, tx, inputAuth); err != nil {
			return err
		}

		response = resp
		return nil
	})
	return response, err
}

func (u *userUsecase) Logout(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	err := u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := u.repo.RevokedRefreshToken(ctx, tx, token); err != nil {
			return err
		}
		return nil
	})
	return err
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *userUsecase) tokenGenerate(input *user.User) (*tokenResponse, error) {
	accessToken, err := u.token.GenerateAccessToken(input)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.token.GenerateRefreshToken(input)
	if err != nil {
		return nil, err
	}

	response := &tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}

func (u *userUsecase) inputAuth(userID int64, token string) *user.Auth {
	return &user.Auth{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(consts.RefreshTokenDuration),
	}
}
