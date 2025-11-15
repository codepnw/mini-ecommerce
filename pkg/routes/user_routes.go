package routes

import (
	"fmt"

	userhandler "github.com/codepnw/mini-ecommerce/internal/user/handler"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
)

func (cfg *routeConfig) UserRoutes() error {
	repo := userrepository.NewUserRepository(cfg.db)
	uc, err := userusecase.NewUserUsecase(&userusecase.UserUsecaseConfig{
		Repo:  repo,
		Token: cfg.token,
		Tx:    cfg.tx,
		DB:    cfg.db,
	})
	if err != nil {
		return fmt.Errorf("user usecase config: %w", err)
	}
	handler := userhandler.NewUserHandler(uc)

	auth := cfg.router.Group("/auth")
	{
		// Public
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		// Private
		auth.POST("/refresh-token", cfg.auth.AuthorizedMiddleware(), handler.RefreshToken)
		auth.POST("/logout", cfg.auth.AuthorizedMiddleware(), handler.Logout)
	}

	return nil
}
