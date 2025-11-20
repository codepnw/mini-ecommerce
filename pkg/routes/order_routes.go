package routes

import (
	"fmt"

	cartrepository "github.com/codepnw/mini-ecommerce/internal/cart/repository"
	orderhandler "github.com/codepnw/mini-ecommerce/internal/order/handler"
	orderrepository "github.com/codepnw/mini-ecommerce/internal/order/repository"
	orderusecase "github.com/codepnw/mini-ecommerce/internal/order/usecase"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	"github.com/codepnw/mini-ecommerce/internal/user"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
)

func (cfg *routeConfig) OrderRoutes() {
	prodRepo := productrepository.NewProductRepository(cfg.db)
	cartRepo := cartrepository.NewCartRepository(cfg.db)
	orderRepo := orderrepository.NewOrderRepository(cfg.db)

	uc := orderusecase.NewOrderUsecase(orderRepo, prodRepo, cartRepo, cfg.tx, cfg.db)
	handler := orderhandler.NewOrderHandler(uc)

	orderID := fmt.Sprintf("/:%s", consts.ParamOrderID)
	r := cfg.router.Group("/orders")
	r.Use(cfg.auth.AuthorizedMiddleware())
	{
		r.POST("/", handler.CreateOrder)
		r.GET(orderID, handler.GetOrderDetail)
		r.GET("/", handler.GetMyOrders)
		r.POST(orderID, handler.CancelOrder)
	}

	// For Admin
	admin := cfg.router.Group("/admin/orders")
	admin.Use(cfg.auth.AuthorizedMiddleware(), cfg.auth.RolesRequired(user.RoleAdmin))
	{
		admin.PATCH(fmt.Sprintf("%s/status", orderID), handler.UpdateOrderStatus)
	}
}
