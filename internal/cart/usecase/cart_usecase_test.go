package cartusecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/codepnw/mini-ecommerce/internal/cart"
	cartrepository "github.com/codepnw/mini-ecommerce/internal/cart/repository"
	cartusecase "github.com/codepnw/mini-ecommerce/internal/cart/usecase"
	"github.com/codepnw/mini-ecommerce/internal/product"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return fn(nil)
}

func TestAddItemToCart(t *testing.T) {
	type testCase struct {
		name        string
		productID   int64
		quantity    int
		mockFn      func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:      "success",
			productID: 101,
			quantity:  3,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockProd := &product.Product{
					ID:    101,
					Stock: 10,
				}
				mockProdRepo.EXPECT().FindByID(gomock.Any(), int64(101)).Return(mockProd, nil).Times(1)

				mockCart := &cart.Cart{
					CartID: "uuid-001",
				}
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)

				mockCartRepo.EXPECT().UpsertItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

				// Return getCartView
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)
				mockItems := []*cartrepository.CartItemDB{
					{CartItemID: "uuid-cart-item-001"},
					{CartItemID: "uuid-cart-item-002"},
					{CartItemID: "uuid-cart-item-003"},
				}
				mockCartRepo.EXPECT().GetCartItems(gomock.Any(), "uuid-001").Return(mockItems, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:      "fail - product not found",
			productID: 101,
			quantity:  3,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockProdRepo.EXPECT().FindByID(gomock.Any(), int64(101)).Return(nil, errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name:      "fail - product not enough",
			productID: 101,
			quantity:  5,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockProd := &product.Product{
					ID:    101,
					Stock: 1,
				}
				mockProdRepo.EXPECT().FindByID(gomock.Any(), int64(101)).Return(mockProd, nil).Times(1)
			},
			expectedErr: errs.ErrProductNotEnough,
		},
		{
			name:      "fail - db get or create cart",
			productID: 101,
			quantity:  5,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockProd := &product.Product{
					ID:    101,
					Stock: 10,
				}
				mockProdRepo.EXPECT().FindByID(gomock.Any(), int64(101)).Return(mockProd, nil).Times(1)

				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCartRepo := cartrepository.NewMockCartRepository(ctrl)
			mockProdRepo := productrepository.NewMockProductRepository(ctrl)
			mockTx := &mockTxManager{}

			uc := cartusecase.NewCartUsecase(mockCartRepo, mockProdRepo, mockTx)

			tc.mockFn(mockCartRepo, mockProdRepo, mockTx)

			result, err := uc.AddItemToCart(context.Background(), tc.productID, tc.quantity)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUpdateItemQuantity(t *testing.T) {
	type testCase struct {
		name        string
		cartItemID  int64
		newQuantity int
		mockFn      func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:        "success",
			cartItemID:  101,
			newQuantity: 3,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockCart := &cart.Cart{CartID: "uuid-001"}
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)

				mockItem := &cart.CartItem{ID: 101, ProductID: 100}
				mockCartRepo.EXPECT().GetCartItemForUpdate(gomock.Any(), gomock.Any(), int64(101), mockCart.CartID).Return(mockItem, nil).Times(1)

				mockProd := &product.Product{ID: 100, Stock: 10}
				mockProdRepo.EXPECT().FindByIDForUpdate(gomock.Any(), gomock.Any(), int64(100)).Return(mockProd, nil).Times(1)

				mockCartRepo.EXPECT().UpdateItemQuantity(gomock.Any(), gomock.Any(), mockCart.CartID, int64(101), 3).Return(nil).Times(1)

				// Return getCartView
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)
				mockItems := []*cartrepository.CartItemDB{
					{CartItemID: "uuid-cart-item-001"},
					{CartItemID: "uuid-cart-item-002"},
					{CartItemID: "uuid-cart-item-003"},
				}
				mockCartRepo.EXPECT().GetCartItems(gomock.Any(), mockCart.CartID).Return(mockItems, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:        "fail - get or create cart",
			cartItemID:  101,
			newQuantity: 3,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
		{
			name:        "fail - item not in cart",
			cartItemID:  101,
			newQuantity: 3,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockCart := &cart.Cart{CartID: "uuid-001"}
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)

				mockCartRepo.EXPECT().GetCartItemForUpdate(gomock.Any(), gomock.Any(), int64(101), mockCart.CartID).Return(nil, errs.ErrItemNotInCart).Times(1)
			},
			expectedErr: errs.ErrItemNotInCart,
		},
		{
			name:        "fail - product not found",
			cartItemID:  101,
			newQuantity: 3,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockCart := &cart.Cart{CartID: "uuid-001"}
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)

				mockItem := &cart.CartItem{ID: 101, ProductID: 100}
				mockCartRepo.EXPECT().GetCartItemForUpdate(gomock.Any(), gomock.Any(), int64(101), mockCart.CartID).Return(mockItem, nil).Times(1)

				mockProdRepo.EXPECT().FindByIDForUpdate(gomock.Any(), gomock.Any(), int64(100)).Return(nil, errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name:        "",
			cartItemID:  101,
			newQuantity: 12,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockCart := &cart.Cart{CartID: "uuid-001"}
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)

				mockItem := &cart.CartItem{ID: 101, ProductID: 100}
				mockCartRepo.EXPECT().GetCartItemForUpdate(gomock.Any(), gomock.Any(), int64(101), mockCart.CartID).Return(mockItem, nil).Times(1)

				mockProd := &product.Product{ID: 100, Stock: 10}
				mockProdRepo.EXPECT().FindByIDForUpdate(gomock.Any(), gomock.Any(), int64(100)).Return(mockProd, nil).Times(1)
			},
			expectedErr: errs.ErrProductNotEnough,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCartRepo := cartrepository.NewMockCartRepository(ctrl)
			mockProdRepo := productrepository.NewMockProductRepository(ctrl)
			mockTx := &mockTxManager{}

			uc := cartusecase.NewCartUsecase(mockCartRepo, mockProdRepo, mockTx)

			tc.mockFn(mockCartRepo, mockProdRepo, mockTx)

			result, err := uc.UpdateItemQuantity(context.Background(), tc.cartItemID, tc.newQuantity)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestRemoveItemFromCart(t *testing.T) {
	type testCase struct {
		name        string
		cartItemID  int64
		mockFn      func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:       "success",
			cartItemID: 100,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockCart := &cart.Cart{CartID: "uuid-001"}
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)

				mockCartRepo.EXPECT().RemoveItem(gomock.Any(), gomock.Any(), mockCart.CartID, int64(100)).Return(nil).Times(1)

				// Return getCartView
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)
				mockItems := []*cartrepository.CartItemDB{
					{CartItemID: "uuid-cart-item-001"},
					{CartItemID: "uuid-cart-item-002"},
					{CartItemID: "uuid-cart-item-003"},
				}
				mockCartRepo.EXPECT().GetCartItems(gomock.Any(), mockCart.CartID).Return(mockItems, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:       "fail - remove item",
			cartItemID: 100,
			mockFn: func(mockCartRepo *cartrepository.MockCartRepository, mockProdRepo *productrepository.MockProductRepository, mockTx *mockTxManager) {
				mockCart := &cart.Cart{CartID: "uuid-001"}
				mockCartRepo.EXPECT().GetOrCreateActiveCart(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCart, nil).Times(1)

				mockCartRepo.EXPECT().RemoveItem(gomock.Any(), gomock.Any(), mockCart.CartID, int64(100)).Return(errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCartRepo := cartrepository.NewMockCartRepository(ctrl)
			mockProdRepo := productrepository.NewMockProductRepository(ctrl)
			mockTx := &mockTxManager{}

			uc := cartusecase.NewCartUsecase(mockCartRepo, mockProdRepo, mockTx)

			tc.mockFn(mockCartRepo, mockProdRepo, mockTx)

			result, err := uc.RemoveItemFromCart(context.Background(), tc.cartItemID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
