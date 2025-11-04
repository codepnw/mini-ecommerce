package productusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/codepnw/mini-ecommerce/internal/product"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	productusecase "github.com/codepnw/mini-ecommerce/internal/product/usecase"
	"github.com/codepnw/mini-ecommerce/internal/user"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUsecase(t *testing.T) {
	type testCase struct {
		name        string
		input       *product.Product
		mockFn      func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success create",
			input: &product.Product{
				Name:    "IPhone 17",
				Price:   35900,
				Stock:   10,
				OwnerID: 10,
				SKU:     "apple-iphone-17",
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockProduct := &product.Product{
					ID:   1,
					Name: "IPhone 17",
					SKU:  "apple-iphone-17",
				}
				mockRepo.EXPECT().SKUExists(gomock.Any(), "apple-iphone-17").Return(false, nil).Times(1)

				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(mockProduct, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail sku aleady exists",
			input: &product.Product{SKU: "apple-iphone-17"},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockRepo.EXPECT().SKUExists(gomock.Any(), "apple-iphone-17").Return(true, nil).Times(1)
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
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockRepo.EXPECT().SKUExists(gomock.Any(), "apple-iphone-17").Return(false, nil).Times(1)

				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Dependencies
			mockRepo := productrepository.NewMockProductRepository(ctrl)
			mockUserUc := userusecase.NewMockUserUsecase(ctrl)

			uc := productusecase.NewProductUsecase(mockRepo, mockUserUc)
			if uc == nil {
				t.Fatalf("usecase is nil")
			}

			tc.mockFn(mockRepo, mockUserUc)

			result, err := uc.Create(context.Background(), tc.input)

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
		input       *product.Product
		mockFn      func(mockRepo *productrepository.MockProductRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success get product",
			input: &product.Product{
				ID: 1,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository) {
				mockProduct := &product.Product{
					ID:   1,
					Name: "IPhone 17",
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail not found",
			input: &product.Product{
				ID: 1,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository) {
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(nil, errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name: "fail get product",
			input: &product.Product{
				ID: 1,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository) {
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := productrepository.NewMockProductRepository(ctrl)
		mockUserUC := userusecase.NewMockUserUsecase(ctrl)
		uc := productusecase.NewProductUsecase(mockRepo, mockUserUC)
		if uc == nil {
			t.Fatalf("product usecase is nil")
		}

		tc.mockFn(mockRepo)

		result, err := uc.GetByID(context.Background(), tc.input.ID)

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
			name:  "success get products",
			input: nil,
			mockFn: func(mockRepo *productrepository.MockProductRepository) {
				mockProducts := []*product.Product{
					{ID: 1, Name: "IPhone"},
					{ID: 2, Name: "Macbook"},
				}
				mockRepo.EXPECT().List(gomock.Any()).Return(mockProducts, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail get products",
			input: &product.Product{
				ID: 1,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository) {
				mockRepo.EXPECT().List(gomock.Any()).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := productrepository.NewMockProductRepository(ctrl)
		mockUserUC := userusecase.NewMockUserUsecase(ctrl)
		uc := productusecase.NewProductUsecase(mockRepo, mockUserUC)
		if uc == nil {
			t.Fatalf("product usecase is nil")
		}

		tc.mockFn(mockRepo)

		result, err := uc.GetAll(context.Background())

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
		mockFn      func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase, input *product.Product)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success update product",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
				Stock:   20,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase, input *product.Product) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockProduct := &product.Product{
					ID:      1,
					OwnerID: 10,
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)

				mockRepo.EXPECT().Update(gomock.Any(), input).Return(input, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail user not found",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase, input *product.Product) {
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(nil, errs.ErrUserNotFound).Times(1)
			},
			expectedErr: errs.ErrUserNotFound,
		},
		{
			name: "fail product not found",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase, input *product.Product) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(nil, errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name: "fail no permission",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase, input *product.Product) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockProduct := &product.Product{
					ID:      1,
					OwnerID: 11,
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)
			},
			expectedErr: errs.ErrNoPermissions,
		},
		{
			name: "fail update product",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase, input *product.Product) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockProduct := &product.Product{
					ID:      1,
					OwnerID: 10,
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)

				mockUpdateProduct := &product.Product{
					ID:    1,
					Stock: 20,
				}
				mockRepo.EXPECT().Update(gomock.Any(), mockUpdateProduct).Return(nil, errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := productrepository.NewMockProductRepository(ctrl)
		mockUserUC := userusecase.NewMockUserUsecase(ctrl)
		uc := productusecase.NewProductUsecase(mockRepo, mockUserUC)
		if uc == nil {
			t.Fatalf("product usecase is nil")
		}

		tc.mockFn(mockRepo, mockUserUC, tc.input)

		result, err := uc.Update(context.Background(), tc.input.OwnerID, tc.input)

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
		mockFn      func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success delete product",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockProduct := &product.Product{
					ID:      1,
					OwnerID: 10,
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)

				mockRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail user not found",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(nil, errs.ErrUserNotFound).Times(1)
			},
			expectedErr: errs.ErrUserNotFound,
		},
		{
			name: "fail product not found",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockProduct := &product.Product{
					ID:      1,
					OwnerID: 10,
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)

				mockRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(errs.ErrProductNotFound).Times(1)
			},
			expectedErr: errs.ErrProductNotFound,
		},
		{
			name: "fail no permission",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockProduct := &product.Product{
					ID:      1,
					OwnerID: 11,
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)
			},
			expectedErr: errs.ErrNoPermissions,
		},
		{
			name: "fail delete product",
			input: &product.Product{
				ID:      1,
				OwnerID: 10,
			},
			mockFn: func(mockRepo *productrepository.MockProductRepository, mockUserUc *userusecase.MockUserUsecase) {
				mockUser := &user.User{
					ID: 10,
				}
				mockUserUc.EXPECT().GetUser(gomock.Any(), int64(10)).Return(mockUser, nil).Times(1)

				mockProduct := &product.Product{
					ID:      1,
					OwnerID: 10,
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockProduct, nil).Times(1)

				mockRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := productrepository.NewMockProductRepository(ctrl)
		mockUserUC := userusecase.NewMockUserUsecase(ctrl)
		uc := productusecase.NewProductUsecase(mockRepo, mockUserUC)
		if uc == nil {
			t.Fatalf("product usecase is nil")
		}

		tc.mockFn(mockRepo, mockUserUC)

		err := uc.Delete(context.Background(), tc.input.OwnerID, tc.input.ID)

		if tc.expectedErr != nil {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
