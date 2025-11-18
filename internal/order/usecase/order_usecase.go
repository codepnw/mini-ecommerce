package orderusecase

import (
	"context"
	"database/sql"
	"time"

	cartrepository "github.com/codepnw/mini-ecommerce/internal/cart/repository"
	"github.com/codepnw/mini-ecommerce/internal/order"
	orderrepository "github.com/codepnw/mini-ecommerce/internal/order/repository"
	"github.com/codepnw/mini-ecommerce/internal/product"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/pkg/auth"
	"github.com/codepnw/mini-ecommerce/pkg/database"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context) (*order.Order, error)
	GetOrderDetail(ctx context.Context, orderID int64) (*OrderView, error)
	GetMyOrders(ctx context.Context) ([]*OrderListView, error)
}

type orderUsecase struct {
	orderRepo   orderrepository.OrderRepository
	productRepo productrepository.ProductRepository
	cartRepo    cartrepository.CartRepository
	tx          database.TxManager
}

func NewOrderUsecase(
	orderRepo orderrepository.OrderRepository,
	productRepo productrepository.ProductRepository,
	cartRepo cartrepository.CartRepository,
	tx database.TxManager,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		tx:          tx,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context) (*order.Order, error) {
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, errs.ErrUnauthorized
	}

	var newOrder *order.Order

	err := u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		cartData, err := u.cartRepo.GetActiveCartByUserID(ctx, tx, userID)
		if err != nil {
			return err
		}

		items, err := u.cartRepo.GetCartItems(ctx, tx, cartData.CartID)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return errs.ErrCartIsEmpty
		}

		lockedProducts := make(map[int64]*product.Product) // Map productID -> *Product
		var totalPrice float64

		// Total Price & Lock Current Product Data
		for _, i := range items {
			product, err := u.productRepo.FindByIDForUpdate(ctx, tx, i.ProductID)
			if err != nil {
				return err
			}
			if product.Stock < i.Quantity {
				return errs.ErrProductNotEnough
			}

			totalPrice += (product.Price * float64(i.Quantity))

			lockedProducts[i.ProductID] = product
		}

		// Create Order
		orderHeader := &order.Order{
			UserID: userID,
			Total:  totalPrice,
			Status: string(order.StatusPending), // Default Status
		}
		newOrderID, err := u.orderRepo.CreateOrder(ctx, tx, orderHeader)
		if err != nil {
			return err
		}

		// New Order
		orderHeader.ID = newOrderID
		newOrder = orderHeader

		// Create Order Items
		for _, i := range items {
			// Lock Product Data
			lockedProduct := lockedProducts[i.ProductID]

			oi := &order.OrderItem{
				OrderID:         newOrderID,
				ProductID:       i.ProductID,
				Quantity:        i.Quantity,
				PriceAtPurchase: lockedProduct.Price, // Current Price
			}
			if err := u.orderRepo.CreateOrderItem(ctx, tx, oi); err != nil {
				return err
			}

			// Decrease Stock
			if err := u.productRepo.DecreaseStock(ctx, tx, i.ProductID, i.Quantity); err != nil {
				return err
			}
		}
		// Clear Cart
		if err := u.cartRepo.ClearCart(ctx, tx, cartData.CartID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (u *orderUsecase) GetOrderDetail(ctx context.Context, orderID int64) (*OrderView, error) {
	// Get UserID
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, errs.ErrUnauthorized
	}

	// Get Order
	orderData, err := u.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if orderData.UserID != userID {
		return nil, errs.ErrNoPermissions
	}

	// Get Items
	itemsData, err := u.orderRepo.GetOrderItems(ctx, orderData.ID)
	if err != nil {
		return nil, err
	}

	// Map Struct -> View
	itemViews := make([]*OrderItemView, 0)
	for _, i := range itemsData {
		itemViews = append(itemViews, &OrderItemView{
			ProductID:       i.ProductID,
			ProductName:     i.ProductName,
			PriceAtPurchase: i.PriceAtPurchase,
			Quantity:        i.Quantity,
			Total:           i.PriceAtPurchase * float64(i.Quantity),
		})
	}

	return &OrderView{
		ID:        orderData.ID,
		Status:    orderData.Status,
		Total:     orderData.Total,
		CreatedAt: orderData.CreatedAt.Format(time.RFC3339),
		Items:     itemViews,
	}, nil
}

func (u *orderUsecase) GetMyOrders(ctx context.Context) ([]*OrderListView, error) {
	// Get UserID
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, errs.ErrUnauthorized
	}

	// Get Orders
	orderData, err := u.orderRepo.GetMyOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Map Struct -> View
	orderView := make([]*OrderListView, 0)
	for _, i := range orderData {
		orderView = append(orderView, &OrderListView{
			ID:        i.ID,
			Status:    i.Status,
			Total:     i.Total,
			CreatedAt: i.CreatedAt.Format(time.RFC3339),
		})
	}
	return orderView, nil
}
