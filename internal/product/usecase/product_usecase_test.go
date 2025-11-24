package productusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/codepnw/mini-ecommerce/internal/product"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/pkg/auth"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUsecase(t *testing.T) {
	type testCase struct {
		name        string
		input       *product.Product
		mockFn      func(mockRepo *productrepository.MockProductRepository, input *product.Product)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			input: &product.Product{
				Name:    "IPhone 17",
				Price:   35900,
				Stock:   10,
				OwnerID: 10,
				SKU:     "apple-iphone-17",
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				mockRepo.EXPECT().SKUExists(gomock.Any(), input.SKU).Return(false, nil).Times(1)

				mockRepo.EXPECT().Insert(gomock.Any(), input).Return(p, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail sku aleady exists",
			input: &product.Product{SKU: "apple-iphone-17"},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				mockRepo.EXPECT().SKUExists(gomock.Any(), input.SKU).Return(true, nil).Times(1)
			},
			expectedErr: errs.ErrProductSKUExists,
		},
		{
			name: "fail create product",
			input: &product.Product{
				Name:    "IPhone 17",
				Price:   35900,
				Stock:   10,
				OwnerID: 10,
				SKU:     "apple-iphone-17",
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				mockRepo.EXPECT().SKUExists(gomock.Any(), input.SKU).Return(false, nil).Times(1)

				mockRepo.EXPECT().Insert(gomock.Any(), input).Return(nil, errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockRepo := setup(t)
			ctx := mockUserClaims()

			tc.mockFn(mockRepo, tc.input)

			// Create Usecase
			result, err := uc.Create(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.ID)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	type testCase struct {
		name        string
		productID   int64
		mockFn      func(mockRepo *productrepository.MockProductRepository, productID int64)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:      "success",
			productID: 1,
			mockFn: func(mockRepo *productrepository.MockProductRepository, productID int64) {
				p := mockProduct()
				mockRepo.EXPECT().FindByID(gomock.Any(), productID).Return(p, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:      "fail not found",
			productID: 100,
			mockFn: func(mockRepo *productrepository.MockProductRepository, productID int64) {
				mockRepo.EXPECT().FindByID(gomock.Any(), productID).Return(nil, errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name:      "fail get product",
			productID: 100,
			mockFn: func(mockRepo *productrepository.MockProductRepository, productID int64) {
				mockRepo.EXPECT().FindByID(gomock.Any(), productID).Return(nil, errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		// Setup
		uc, mockRepo := setup(t)

		tc.mockFn(mockRepo, tc.productID)

		// GetByID Usecase
		result, err := uc.GetByID(context.Background(), tc.productID)

		if tc.expectedErr != nil {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
			assert.Nil(t, result)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.ID)
		}
	}
}

func TestGetProducts(t *testing.T) {
	type testCase struct {
		name        string
		input       *product.Product
		mockFn      func(mockRepo *productrepository.MockProductRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success",
			input: nil,
			mockFn: func(mockRepo *productrepository.MockProductRepository) {
				p := mockProduct()
				list := []*product.Product{
					{ID: p.ID},
					{ID: p.ID + 1},
				}
				mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(list, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail get products",
			input: nil,
			mockFn: func(mockRepo *productrepository.MockProductRepository) {
				mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		// Setup
		uc, mockRepo := setup(t)

		tc.mockFn(mockRepo)

		// List Usecase
		result, err := uc.List(context.Background(), &product.ProductFilter{
			Limit: 20,
		})

		if tc.expectedErr != nil {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
			assert.Nil(t, result)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, result)
		}
	}
}

func TestUpdateProduct(t *testing.T) {
	type testCase struct {
		name        string
		input       *product.Product
		mockFn      func(mockRepo *productrepository.MockProductRepository, input *product.Product)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
				Stock:   20,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				mockRepo.EXPECT().FindByID(gomock.Any(), input.ID).Return(p, nil).Times(1)

				mockRepo.EXPECT().Update(gomock.Any(), input).Return(input, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail product not found",
			input: &product.Product{
				ID: 1,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				mockRepo.EXPECT().FindByID(gomock.Any(), input.ID).Return(nil, errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name: "fail no permission",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				p.OwnerID = 11
				mockRepo.EXPECT().FindByID(gomock.Any(), input.ID).Return(p, nil).Times(1)
			},
			expectedErr: errs.ErrNoPermissions,
		},
		{
			name: "fail update product",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				mockRepo.EXPECT().FindByID(gomock.Any(), input.ID).Return(p, nil).Times(1)

				mockRepo.EXPECT().Update(gomock.Any(), input).Return(nil, errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		// Setup
		uc, mockRepo := setup(t)
		ctx := mockUserClaims()

		tc.mockFn(mockRepo, tc.input)

		// Update Usecase
		result, err := uc.Update(ctx, tc.input)

		if tc.expectedErr != nil {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
			assert.Nil(t, result)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, result)
		}
	}
}

func TestDeleteProduct(t *testing.T) {
	type testCase struct {
		name        string
		input       *product.Product
		mockFn      func(mockRepo *productrepository.MockProductRepository, input *product.Product)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				mockRepo.EXPECT().FindByID(gomock.Any(), input.ID).Return(p, nil).Times(1)

				mockRepo.EXPECT().Delete(gomock.Any(), input.ID).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail product not found",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				mockRepo.EXPECT().FindByID(gomock.Any(), input.ID).Return(p, nil).Times(1)

				mockRepo.EXPECT().Delete(gomock.Any(), input.ID).Return(errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name: "fail no permission",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				p.OwnerID = 11
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(p, nil).Times(1)
			},
			expectedErr: errs.ErrNoPermissions,
		},
		{
			name: "fail delete product",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, input *product.Product) {
				p := mockProduct()
				mockRepo.EXPECT().FindByID(gomock.Any(), input.ID).Return(p, nil).Times(1)

				mockRepo.EXPECT().Delete(gomock.Any(), input.ID).Return(errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		uc, mockRepo := setup(t)
		ctx := mockUserClaims()

		tc.mockFn(mockRepo, tc.input)

		err := uc.Delete(ctx, tc.input.ID)

		if tc.expectedErr != nil {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

// ============ Helper ================
// ------------------------------------
func setup(t *testing.T) (productusecase.ProductUsecase, *productrepository.MockProductRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := productrepository.NewMockProductRepository(ctrl)
	uc := productusecase.NewProductUsecase(mockRepo)

	return uc, mockRepo
}

func mockUserClaims() context.Context {
	mockUser := &jwt.UserClaims{
		ID:    10,
		Email: "example@mail.com",
		Role:  "user",
	}
	return auth.SetCurrentUser(context.Background(), mockUser)
}

func mockProduct() *product.Product {
	return &product.Product{
		ID:      100,
		OwnerID: 10,
		Name:    "Macbook",
		Stock:   20,
		SKU:     "mock-product",
	}
}

var errDBMock = errors.New("db error")
