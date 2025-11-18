package orderusecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/codepnw/mini-ecommerce/internal/cart"
	cartrepository "github.com/codepnw/mini-ecommerce/internal/cart/repository"
	"github.com/codepnw/mini-ecommerce/internal/order"
	orderrepository "github.com/codepnw/mini-ecommerce/internal/order/repository"
	orderusecase "github.com/codepnw/mini-ecommerce/internal/order/usecase"
	"github.com/codepnw/mini-ecommerce/internal/product"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return fn(nil)
}

type mockDB struct{}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return nil
}

func (m *mockDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, nil
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, nil
}

func TestCreateOrder(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		mockFn      func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:   "success",
			userID: 10,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockCart := &cart.Cart{CartID: "cart-001", UserID: sql.NullInt64{Int64: 10}}
				cartRepo.EXPECT().GetActiveCartByUserID(gomock.Any(), gomock.Any(), int64(10)).Return(mockCart, nil).Times(1)

				mockItems := []*cartrepository.CartItemDB{
					{CartItemID: "item-001", ProductID: 1, Price: 100, Quantity: 2},
					{CartItemID: "item-002", ProductID: 2, Price: 80, Quantity: 2},
				}
				cartRepo.EXPECT().GetCartItems(gomock.Any(), gomock.Any(), mockCart.CartID).Return(mockItems, nil).Times(1)

				var expectedTotal float64

				for _, i := range mockItems {
					mockProduct := &product.Product{
						ID:    i.ProductID,
						Price: i.Price,
						Stock: 100,
					}
					prodRepo.EXPECT().FindByIDForUpdate(gomock.Any(), gomock.Any(), i.ProductID).Return(mockProduct, nil).Times(1)
					expectedTotal += (i.Price * float64(i.Quantity))
				}

				var mockOrderID int64 = 1
				mockOrderHeader := &order.Order{
					UserID: mockCart.UserID.Int64,
					Total:  expectedTotal,
					Status: string(order.StatusPending),
				}
				orderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), mockOrderHeader).Return(mockOrderID, nil).Times(1)

				// Create Order Items
				for _, i := range mockItems {
					mockOI := &order.OrderItem{
						OrderID:         mockOrderID,
						ProductID:       i.ProductID,
						Quantity:        i.Quantity,
						PriceAtPurchase: i.Price,
					}
					orderRepo.EXPECT().CreateOrderItem(gomock.Any(), gomock.Any(), mockOI).Return(nil).Times(1)

					prodRepo.EXPECT().DecreaseStock(gomock.Any(), gomock.Any(), i.ProductID, i.Quantity).Return(nil).Times(1)
				}
				cartRepo.EXPECT().ClearCart(gomock.Any(), gomock.Any(), mockCart.CartID).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:   "fail product not enough",
			userID: 10,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockCart := &cart.Cart{CartID: "cart-001", UserID: sql.NullInt64{Int64: 10}}
				cartRepo.EXPECT().GetActiveCartByUserID(gomock.Any(), gomock.Any(), int64(10)).Return(mockCart, nil).Times(1)

				mockItems := []*cartrepository.CartItemDB{
					{CartItemID: "item-001", ProductID: 1, Price: 100, Quantity: 2},
				}
				cartRepo.EXPECT().GetCartItems(gomock.Any(), gomock.Any(), mockCart.CartID).Return(mockItems, nil).Times(1)

				for _, i := range mockItems {
					mockProduct := &product.Product{
						ID:    i.ProductID,
						Price: i.Price,
						Stock: 1,
					}
					prodRepo.EXPECT().FindByIDForUpdate(gomock.Any(), gomock.Any(), i.ProductID).Return(mockProduct, nil).Times(1)
				}
			},
			expectedErr: errs.ErrProductNotEnough,
		},
		{
			name:   "fail db error",
			userID: 10,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockCart := &cart.Cart{CartID: "cart-001", UserID: sql.NullInt64{Int64: 10}}
				cartRepo.EXPECT().GetActiveCartByUserID(gomock.Any(), gomock.Any(), int64(10)).Return(mockCart, nil).Times(1)

				mockItems := []*cartrepository.CartItemDB{
					{CartItemID: "item-001", ProductID: 1, Price: 100, Quantity: 2},
					{CartItemID: "item-002", ProductID: 2, Price: 80, Quantity: 2},
				}
				cartRepo.EXPECT().GetCartItems(gomock.Any(), gomock.Any(), mockCart.CartID).Return(mockItems, nil).Times(1)

				var expectedTotal float64

				for _, i := range mockItems {
					mockProduct := &product.Product{
						ID:    i.ProductID,
						Price: i.Price,
						Stock: 100,
					}
					prodRepo.EXPECT().FindByIDForUpdate(gomock.Any(), gomock.Any(), i.ProductID).Return(mockProduct, nil).Times(1)
					expectedTotal += (i.Price * float64(i.Quantity))
				}

				var mockOrderID int64 = 1
				mockOrderHeader := &order.Order{
					UserID: mockCart.UserID.Int64,
					Total:  expectedTotal,
					Status: string(order.StatusPending),
				}
				orderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), mockOrderHeader).Return(mockOrderID, nil).Times(1)

				// Create Order Items
				for _, i := range mockItems {
					mockOI := &order.OrderItem{
						OrderID:         mockOrderID,
						ProductID:       i.ProductID,
						Quantity:        i.Quantity,
						PriceAtPurchase: i.Price,
					}
					orderRepo.EXPECT().CreateOrderItem(gomock.Any(), gomock.Any(), mockOI).Return(nil).Times(1)

					prodRepo.EXPECT().DecreaseStock(gomock.Any(), gomock.Any(), i.ProductID, i.Quantity).Return(nil).Times(1)
				}
				cartRepo.EXPECT().ClearCart(gomock.Any(), gomock.Any(), mockCart.CartID).Return(errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := orderrepository.NewMockOrderRepository(ctrl)
			cartRepo := cartrepository.NewMockCartRepository(ctrl)
			prodRepo := productrepository.NewMockProductRepository(ctrl)
			mockTx := &mockTxManager{}
			mockDB := &mockDB{}

			uc := orderusecase.NewOrderUsecase(orderRepo, prodRepo, cartRepo, mockTx, mockDB)

			tc.mockFn(orderRepo, prodRepo, cartRepo)

			// Set UserID
			ctx := context.Background()
			if tc.userID != 0 {
				ctx = auth.SetUserID(ctx, tc.userID)
			}

			result, err := uc.CreateOrder(ctx)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestGetOrderDetail(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		orderID     int64
		mockFn      func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:    "success",
			userID:  10,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockOrder := &order.Order{
					ID:     100,
					UserID: 10,
					Status: "pending",
				}
				orderRepo.EXPECT().GetOrder(gomock.Any(), int64(100)).Return(mockOrder, nil).Times(1)

				mockItems := []*orderrepository.OrderItemDetail{
					{
						ProductID:       100,
						PriceAtPurchase: 35000.00,
						ProductName:     "macbook",
						Quantity:        2,
					},
					{
						ProductID:       101,
						PriceAtPurchase: 25000.00,
						ProductName:     "ipad",
						Quantity:        1,
					},
				}
				orderRepo.EXPECT().GetOrderItems(gomock.Any(), gomock.Any(), mockOrder.ID).Return(mockItems, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:    "fail unauthorized",
			userID:  0,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
			},
			expectedErr: errs.ErrNoPermissions,
		},
		{
			name:    "fail get order",
			userID:  10,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				orderRepo.EXPECT().GetOrder(gomock.Any(), int64(100)).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
		{
			name:    "fail no permissions",
			userID:  10,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockOrder := &order.Order{
					ID:     100,
					UserID: 11,
					Status: "pending",
				}
				orderRepo.EXPECT().GetOrder(gomock.Any(), int64(100)).Return(mockOrder, nil).Times(1)
			},
			expectedErr: errs.ErrNoPermissions,
		},
		{
			name:    "fail get items",
			userID:  10,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockOrder := &order.Order{
					ID:     100,
					UserID: 10,
					Status: "pending",
				}
				orderRepo.EXPECT().GetOrder(gomock.Any(), int64(100)).Return(mockOrder, nil).Times(1)

				orderRepo.EXPECT().GetOrderItems(gomock.Any(), gomock.Any(), mockOrder.ID).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := orderrepository.NewMockOrderRepository(ctrl)
		cartRepo := cartrepository.NewMockCartRepository(ctrl)
		prodRepo := productrepository.NewMockProductRepository(ctrl)
		mockTx := &mockTxManager{}
		mockDB := &mockDB{}

		uc := orderusecase.NewOrderUsecase(orderRepo, prodRepo, cartRepo, mockTx, mockDB)

		tc.mockFn(orderRepo, prodRepo, cartRepo)

		// Set UserID
		ctx := context.Background()
		if tc.userID != 0 {
			ctx = auth.SetUserID(ctx, tc.userID)
		}

		result, err := uc.GetOrderDetail(ctx, tc.orderID)

		if tc.expectedErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, result)
		}
	}
}

func TestGetMyOrders(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		mockFn      func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:   "success",
			userID: 10,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockOrders := []*order.Order{
					{ID: 1, Status: "pending", Total: 100, CreatedAt: time.Now()},
					{ID: 2, Status: "pending", Total: 200, CreatedAt: time.Now()},
				}
				orderRepo.EXPECT().GetMyOrders(gomock.Any(), int64(10)).Return(mockOrders, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:   "fail unauthorized",
			userID: 0,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
			},
			expectedErr: errs.ErrUnauthorized,
		},
		{
			name:   "fail get orders",
			userID: 10,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				orderRepo.EXPECT().GetMyOrders(gomock.Any(), int64(10)).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := orderrepository.NewMockOrderRepository(ctrl)
			cartRepo := cartrepository.NewMockCartRepository(ctrl)
			prodRepo := productrepository.NewMockProductRepository(ctrl)
			mockTx := &mockTxManager{}
			mockDB := &mockDB{}

			uc := orderusecase.NewOrderUsecase(orderRepo, prodRepo, cartRepo, mockTx, mockDB)

			tc.mockFn(orderRepo, prodRepo, cartRepo)

			// Set UserID
			ctx := context.Background()
			if tc.userID != 0 {
				ctx = auth.SetUserID(ctx, tc.userID)
			}

			result, err := uc.GetMyOrders(ctx)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestCancelOrder(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		orderID     int64
		mockFn      func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:    "success",
			userID:  10,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockOrder := &order.Order{
					ID:     100,
					UserID: 10,
					Status: string(order.StatusPending),
				}
				orderRepo.EXPECT().GetOrder(gomock.Any(), int64(100)).Return(mockOrder, nil).Times(1)

				orderRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), int64(mockOrder.ID), string(order.StatusCancelled)).Return(nil).Times(1)

				mockItems := []*orderrepository.OrderItemDetail{
					{ID: 1, ProductID: 100, Quantity: 2},
					{ID: 2, ProductID: 101, Quantity: 5},
				}
				orderRepo.EXPECT().GetOrderItems(gomock.Any(), gomock.Any(), mockOrder.ID).Return(mockItems, nil).Times(1)

				for _, i := range mockItems {
					prodRepo.EXPECT().IncreaseStock(gomock.Any(), gomock.Any(), i.ProductID, i.Quantity).Return(nil).Times(1)
				}
			},
			expectedErr: nil,
		},
		{
			name:    "fail user id",
			userID:  0,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
			},
			expectedErr: errs.ErrUnauthorized,
		},
		{
			name:    "fail cannot cancel",
			userID:  10,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockOrder := &order.Order{
					ID:     100,
					UserID: 10,
					Status: string(order.StatusCancelled),
				}
				orderRepo.EXPECT().GetOrder(gomock.Any(), int64(100)).Return(mockOrder, nil).Times(1)
			},
			expectedErr: errs.ErrCannotCancelOrder,
		},
		{
			name:    "fail get items",
			userID:  10,
			orderID: 100,
			mockFn: func(orderRepo *orderrepository.MockOrderRepository, prodRepo *productrepository.MockProductRepository, cartRepo *cartrepository.MockCartRepository) {
				mockOrder := &order.Order{
					ID:     100,
					UserID: 10,
					Status: string(order.StatusPending),
				}
				orderRepo.EXPECT().GetOrder(gomock.Any(), int64(100)).Return(mockOrder, nil).Times(1)

				orderRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), int64(mockOrder.ID), string(order.StatusCancelled)).Return(nil).Times(1)

				orderRepo.EXPECT().GetOrderItems(gomock.Any(), gomock.Any(), mockOrder.ID).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			orderRepo := orderrepository.NewMockOrderRepository(ctrl)
			cartRepo := cartrepository.NewMockCartRepository(ctrl)
			prodRepo := productrepository.NewMockProductRepository(ctrl)
			mockTx := &mockTxManager{}
			mockDB := &mockDB{}

			uc := orderusecase.NewOrderUsecase(orderRepo, prodRepo, cartRepo, mockTx, mockDB)

			tc.mockFn(orderRepo, prodRepo, cartRepo)

			// Set UserID
			ctx := context.Background()
			if tc.userID != 0 {
				ctx = auth.SetUserID(ctx, tc.userID)
			}

			err := uc.CancelOrder(ctx, tc.orderID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
