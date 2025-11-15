package routes

import (
	"fmt"

	producthandler "github.com/codepnw/mini-ecommerce/internal/product/handler"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
	"github.com/codepnw/mini-ecommerce/internal/user"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
)

func (cfg *routeConfig) ProductRoutes() error {
	userRepo := userrepository.NewUserRepository(cfg.db)
	userUc, err := userusecase.NewUserUsecase(&userusecase.UserUsecaseConfig{
		Repo:  userRepo,
		Token: cfg.token,
		Tx:    cfg.tx,
		DB:    cfg.db,
	})
	if err != nil {
		return err
	}

	repo := productrepository.NewProductRepository(cfg.db)
	uc := productusecase.NewProductUsecase(repo, userUc)
	handler := producthandler.NewProductHandler(uc)

	paramID := fmt.Sprintf("/:%s", consts.ParamProductID)
	public := cfg.router.Group("/products")
	private := cfg.router.Group(
		"/products",
		cfg.auth.AuthorizedMiddleware(),
		cfg.auth.RolesRequired(user.RoleAdmin, user.RoleSeller),
	)
	{
		// Public
		public.GET("/", handler.GetAll)
		public.GET(paramID, handler.GetByID)
		// Admin & Seller
		private.POST("/", handler.Create)
		private.PATCH(paramID, handler.Update)
		private.DELETE(paramID, handler.Delete)
	}

	return nil
}
