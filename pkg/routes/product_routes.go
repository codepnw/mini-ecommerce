package routes

import (
	"fmt"

	producthandler "github.com/codepnw/mini-ecommerce/internal/product/handler"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
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
	r := cfg.router.Group("/products")
	{
		r.POST("/", handler.Create)
		r.GET("/", handler.GetAll)
		r.GET(paramID, handler.GetByID)
		r.PATCH(paramID, handler.Update)
		r.DELETE(paramID, handler.Delete)
	}

	return nil
}
