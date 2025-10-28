package routes

import (
	userhandler "github.com/codepnw/mini-ecommerce/internal/user/handler"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
)

func (cfg *routeConfig) UserRoutes() {
	repo := userrepository.NewUserRepository(cfg.db)
	uc := userusecase.NewUserUsecase(repo, cfg.token)
	handler := userhandler.NewUserHandler(uc)

	r := cfg.router.Group("/users")

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
}
