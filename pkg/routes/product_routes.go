package routes

import (
	"fmt"

	producthandler "github.com/codepnw/mini-ecommerce/internal/product/handler"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
)

func (cfg *routeConfig) ProductRoutes() {
	repo := productrepository.NewProductRepository(cfg.db)
	uc := productusecase.NewProductUsecase(repo)
	handler := producthandler.NewProductHandler(uc)

	paramID := fmt.Sprintf("/:%s", consts.ParamProductID)
	r := cfg.router.Group("/products")
	{
		r.POST("/", handler.Create)
		r.GET("/", handler.GetAll)
		r.GET(paramID, handler.GetByID)
		r.PATCH(paramID, handler.Update)
		r.DELETE(paramID, handler.Delete)
	}
}
