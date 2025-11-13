package routes

import (
	"fmt"

	carthandler "github.com/codepnw/mini-ecommerce/internal/cart/handler"
	cartrepository "github.com/codepnw/mini-ecommerce/internal/cart/repository"
	cartusecase "github.com/codepnw/mini-ecommerce/internal/cart/usecase"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
)

func (cfg *routeConfig) CartRoutes() {
	prodRepo := productrepository.NewProductRepository(cfg.db)
	cartRepo := cartrepository.NewCartRepository(cfg.db)
	uc := cartusecase.NewCartUsecase(cartRepo, prodRepo, cfg.tx)
	handler := carthandler.NewCartHandler(uc)

	var cartItemID = fmt.Sprintf("/items/:%s", consts.CartItemID)
	cartRoutes := cfg.router.Group("/cart")
	{
		cartRoutes.POST("/", handler.AddItemToCart)
		cartRoutes.GET("/", handler.GetCart)
		cartRoutes.PATCH(cartItemID, handler.UpdateItemQuantity)
		cartRoutes.DELETE(cartItemID, handler.RemoveItemFromCart)
	}
}
