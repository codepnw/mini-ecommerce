package userusecase

import (
	"context"

	"github.com/codepnw/mini-ecommerce/internal/consts"
	"github.com/codepnw/mini-ecommerce/internal/user"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/codepnw/mini-ecommerce/pkg/password"
)

type UserUsecase interface {
	Create(ctx context.Context, input *user.User) (*tokenResponse, error)
	Login(ctx context.Context, input *user.User) (*tokenResponse, error)
}

type userUsecase struct {
	repo  userrepository.UserRepository
	token *jwt.JWTToken
}

func NewUserUsecase(repo userrepository.UserRepository, token *jwt.JWTToken) UserUsecase {
	return &userUsecase{
		repo:  repo,
		token: token,
	}
}

func (u *userUsecase) Create(ctx context.Context, input *user.User) (*tokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	hashed, err := password.HashedPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.Password = hashed

	user, err := u.repo.Insert(ctx, input)
	if err != nil {
		return nil, err
	}

	response, err := u.tokenGenerate(user)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *userUsecase) Login(ctx context.Context, input *user.User) (*tokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	user, err := u.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	err = password.ComparePassword(user.Password, input.Password)
	if err != nil {
		return nil, err
	}

	response, err := u.tokenGenerate(user)
	if err != nil {
		return nil, err
	}
	return response, nil
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
