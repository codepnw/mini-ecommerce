package routes

import (
	"fmt"

	producthandler "github.com/codepnw/mini-ecommerce/internal/product/handler"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
	"github.com/codepnw/mini-ecommerce/internal/user"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
)

func (cfg *routeConfig) ProductRoutes() {
	repo := productrepository.NewProductRepository(cfg.db)
	uc := productusecase.NewProductUsecase(repo)
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
}
