package cartusecase

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/codepnw/mini-ecommerce/internal/cart"
	cartrepository "github.com/codepnw/mini-ecommerce/internal/cart/repository"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/pkg/auth"
	"github.com/codepnw/mini-ecommerce/pkg/database"
)

type CartUsecase interface {
	AddItemToCart(ctx context.Context, productID int64, quantity int) (*CartView, error)
	GetCart(ctx context.Context) (*CartView, error)
	UpdateItemQuantity(ctx context.Context, cartItemID int64, newQuantity int) (*CartView, error)
	RemoveItemFromCart(ctx context.Context, cartItemID int64) (*CartView, error)
}

type cartUsecase struct {
	cartRepo    cartrepository.CartRepository
	productRepo productrepository.ProductRepository
	tx          database.TxManager
}

func NewCartUsecase(cartRepo cartrepository.CartRepository, productRepo productrepository.ProductRepository, tx database.TxManager) CartUsecase {
	return &cartUsecase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		tx:          tx,
	}
}

func (u *cartUsecase) AddItemToCart(ctx context.Context, productID int64, quantity int) (*CartView, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Check Product Stock
	product, err := u.productRepo.FindByID(ctx, productID)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			return nil, err
		}
		return nil, err
	}
	if product.Stock < quantity {
		return nil, errs.ErrProductNotEnough
	}

	userID := auth.GetUserID(ctx)
	sessionID := auth.GetSessionID(ctx)
	nullUserID := sql.NullInt64{Int64: userID, Valid: userID > 0}
	nullSessionID := sql.NullString{String: sessionID, Valid: sessionID != ""}

	// Get or Create Cart
	cartData, err := u.cartRepo.GetOrCreateActiveCart(ctx, nullUserID, nullSessionID)
	if err != nil {
		return nil, err
	}

	// Upsert Item
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		item := &cart.CartItem{
			CartID:     cartData.CartID,
			ProductID:  productID,
			Quantity:   quantity,
			PriceAtAdd: product.Price,
		}
		return u.cartRepo.UpsertItem(ctx, tx, item)
	})
	if err != nil {
		return nil, err
	}

	// TODO: Clear Cache (Redis)
	return u.getCartView(ctx)
}

type CartItemView struct {
	CartItemID string `json:"cart_item_id"`
	ProductID  int64  `json:"product_id"`
	Quantity   int    `json:"quantity"`
	// From Products
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	SKU   string  `json:"sku"`
	// Validation
	PriceAtAdd     float64 `json:"-"`
	IsPriceChanged bool    `json:"is_price_changed"`
	IsOutOfStock   bool    `json:"is_out_of_stock"`
}

type CartView struct {
	CartID     string          `json:"cart_id"`
	UserID     *int64          `json:"user_id"`
	Items      []*CartItemView `json:"items"`
	TotalPrice float64         `json:"total_price"`
	TotalItems int             `json:"total_items"`
	HasChanged bool            `json:"has_changed"`
}

func (u *cartUsecase) GetCart(ctx context.Context) (*CartView, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// TODO: Create Cache
	return u.getCartView(ctx)
}

func (u *cartUsecase) getCartView(ctx context.Context) (*CartView, error) {
	userID := auth.GetUserID(ctx)
	sessionID := auth.GetSessionID(ctx)
	nullUserID := sql.NullInt64{Int64: userID, Valid: userID > 0}
	nullSessionID := sql.NullString{String: sessionID, Valid: sessionID != ""}

	// Get or Create Cart
	cartData, err := u.cartRepo.GetOrCreateActiveCart(ctx, nullUserID, nullSessionID)
	if err != nil {
		return nil, err
	}

	// Get Items
	items, err := u.cartRepo.GetCartItems(ctx, cartData.CartID)
	if err != nil {
		return nil, err
	}

	finalItems := make([]*CartItemView, 0)
	var (
		totalPrice float64 = 0
		totalItems         = 0
		hasChanged         = false
	)
	for _, item := range items {
		// Check Stock
		isOutOfStock := item.Quantity > item.Stock
		// Check Current Price
		isPriceChanged := math.Abs(item.PriceAtAdd-item.Price) > 0

		if isOutOfStock || isPriceChanged {
			hasChanged = true
		}

		// Check SKU
		var finalSKU string
		if item.SKU.Valid {
			finalSKU = item.SKU.String
		}

		viewItem := &CartItemView{
			CartItemID:     item.CartItemID,
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			Name:           item.Name,
			Price:          item.Price,
			Stock:          item.Stock,
			SKU:            finalSKU,
			PriceAtAdd:     item.PriceAtAdd,
			IsPriceChanged: isPriceChanged,
			IsOutOfStock:   isOutOfStock,
		}
		finalItems = append(finalItems, viewItem)

		if !isOutOfStock {
			totalPrice += (item.Price * float64(item.Quantity))
		}
		totalItems += item.Quantity
	}

	var finalUserID *int64
	if userID > 0 {
		finalUserID = &userID
	}

	return &CartView{
		CartID:     cartData.CartID,
		UserID:     finalUserID,
		Items:      finalItems,
		TotalPrice: totalPrice,
		TotalItems: totalItems,
		HasChanged: hasChanged,
	}, nil
}

func (u *cartUsecase) UpdateItemQuantity(ctx context.Context, cartItemID int64, newQuantity int) (*CartView, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if newQuantity <= 0 {
		return nil, errs.ErrInvalidQuantity
	}

	userID := auth.GetUserID(ctx)
	sessionID := auth.GetSessionID(ctx)
	nullUserID := sql.NullInt64{Int64: userID, Valid: userID > 0}
	nullSessionID := sql.NullString{String: sessionID, Valid: sessionID != ""}

	// Get or Create Cart
	cartData, err := u.cartRepo.GetOrCreateActiveCart(ctx, nullUserID, nullSessionID)
	if err != nil {
		return nil, err
	}

	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Item Detail
		item, err := u.cartRepo.GetCartItemForUpdate(ctx, tx, cartItemID, cartData.CartID)
		if err != nil {
			return errs.ErrItemNotInCart
		}

		product, err := u.productRepo.FindByIDForUpdate(ctx, tx, item.ProductID)
		if err != nil {
			if errors.Is(err, errs.ErrProductNotFound) {
				return err
			}
			return err
		}
		if product.Stock < newQuantity {
			return errs.ErrProductNotEnough
		}

		return u.cartRepo.UpdateItemQuantity(ctx, tx, cartData.CartID, cartItemID, newQuantity)
	})
	if err != nil {
		return nil, err
	}

	// TODO: Clear Cache
	return u.getCartView(ctx)
}

func (u *cartUsecase) RemoveItemFromCart(ctx context.Context, cartItemID int64) (*CartView, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userID := auth.GetUserID(ctx)
	sessionID := auth.GetSessionID(ctx)
	nullUserID := sql.NullInt64{Int64: userID, Valid: userID > 0}
	nullSessionID := sql.NullString{String: sessionID, Valid: sessionID != ""}

	cartData, err := u.cartRepo.GetOrCreateActiveCart(ctx, nullUserID, nullSessionID)
	if err != nil {
		return nil, err
	}

	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return u.cartRepo.RemoveItem(ctx, tx, cartData.CartID, cartItemID)
	})
	if err != nil {
		return nil, err
	}

	return u.getCartView(ctx)
}
