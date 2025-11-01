package userusecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/codepnw/mini-ecommerce/internal/errs"
	"github.com/codepnw/mini-ecommerce/internal/user"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
	"github.com/codepnw/mini-ecommerce/pkg/config"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/codepnw/mini-ecommerce/pkg/password"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return fn(nil)
}

func TestRegisterUsecse(t *testing.T) {
	type testCase struct {
		name        string
		input       *user.User
		mockFn      func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success user created",
			input: &user.User{Email: "user@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager) {
				mockUser := &user.User{
					ID:    1,
					Email: "user@example.com",
				}
				mockRepo.EXPECT().Insert(gomock.Any(), nil, gomock.Any()).Return(mockUser, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail email already exists",
			input: &user.User{Email: "user@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager) {
				mockRepo.EXPECT().Insert(gomock.Any(), nil, gomock.Any()).Return(nil, errs.ErrEmailAlreadyExists).Times(1)
			},
			expectedErr: errs.ErrEmailAlreadyExists,
		},
		{
			name:  "fail save token",
			input: &user.User{Email: "user@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager) {
				mockUser := &user.User{
					ID:    1,
					Email: "user@example.com",
				}
				mockRepo.EXPECT().Insert(gomock.Any(), nil, gomock.Any()).Return(mockUser, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Dependencies
			mockRepo := userrepository.NewMockUserRepository(ctrl)
			mockTx := &mockTxManager{}
			mockToken, err := jwt.InitJWT(config.JWTConfig{
				SecretKey:  "mock_secret_key",
				RefreshKey: "mock_refresh_key",
			})
			if err != nil {
				t.Fatalf("InitJWT failed: %v", err)
			}

			// NewUsecase
			uc, err := userusecase.NewUserUsecase(&userusecase.UserUsecaseConfig{
				Repo:  mockRepo,
				Tx:    mockTx,
				Token: mockToken,
				DB:    nil,
			})
			if err != nil {
				t.Fatalf("NewUserUsecase failed: %v", err)
			}

			// Mock Repo Response
			tc.mockFn(mockRepo, mockTx)

			// Act
			result, err := uc.Register(context.Background(), tc.input)

			// Assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestLoginUsecase(t *testing.T) {
	type testCase struct {
		name        string
		input       *user.User
		mockFn      func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "fail wrong password",
			input: &user.User{Email: "user@example.com", Password: "wrong_password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager) {
				hashedPassword, _ := password.HashedPassword("correct_password")
				mockUser := &user.User{
					ID:       1,
					Email:    "user@example.com",
					Password: hashedPassword,
				}

				mockRepo.EXPECT().FindByEmail(gomock.Any(), "user@example.com").Return(mockUser, nil).Times(1)
			},
			expectedErr: errs.ErrUserCredentials,
		},
		{
			name:  "success ok",
			input: &user.User{Email: "user@example.com", Password: "correct_password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager) {
				hashedPassword, _ := password.HashedPassword("correct_password")
				mockUser := &user.User{
					ID:       1,
					Email:    "user@example.com",
					Password: hashedPassword,
				}

				mockRepo.EXPECT().FindByEmail(gomock.Any(), "user@example.com").Return(mockUser, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail save token",
			input: &user.User{Email: "user@example.com", Password: "correct_password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager) {
				hashedPassword, _ := password.HashedPassword("correct_password")
				mockUser := &user.User{
					ID:       1,
					Email:    "user@example.com",
					Password: hashedPassword,
				}

				mockRepo.EXPECT().FindByEmail(gomock.Any(), "user@example.com").Return(mockUser, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(errors.New("DB connection error")).Times(1)
			},
			expectedErr: errors.New("DB connection error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Dependencies
			mockRepo := userrepository.NewMockUserRepository(ctrl)
			mockTx := &mockTxManager{}
			mockToken, err := jwt.InitJWT(config.JWTConfig{
				SecretKey:  "mock_secret_key",
				RefreshKey: "mock_refresh_key",
			})
			if err != nil {
				t.Fatalf("InitJWT failed: %v", err)
			}

			// New User Usecase
			uc, err := userusecase.NewUserUsecase(&userusecase.UserUsecaseConfig{
				Repo:  mockRepo,
				Token: mockToken,
				Tx:    mockTx,
				DB:    nil,
			})
			if err != nil {
				t.Fatalf("NewUserUsecase failed: %v", err)
			}

			// Mock Repo Response
			tc.mockFn(mockRepo, mockTx)

			// Act
			result, err := uc.Login(context.Background(), tc.input)

			// Assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestRefreshTokenUsecase(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		mockFn      func(mockRepo *userrepository.MockUserRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success refresh token",
			input: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository) {
				var userID int64 = 1
				mockToken := "mock_refresh_token"
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), mockToken).Return(userID, nil).Times(1)

				mockUser := &user.User{
					ID:    1,
					Email: "user@example.com",
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(mockUser, nil).Times(1)

				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), mockToken).Return(nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail user not found",
			input: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository) {
				var userID int64 = 1
				mockToken := "mock_refresh_token"
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), mockToken).Return(userID, nil).Times(1)

				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, errs.ErrUserNotFound).Times(1)
			},
			expectedErr: errs.ErrUserNotFound,
		},
		{
			name:  "fail save token",
			input: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository) {
				var userID int64 = 1
				mockToken := "mock_refresh_token"
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), mockToken).Return(userID, nil).Times(1)

				mockUser := &user.User{
					ID:    1,
					Email: "user@example.com",
				}
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(mockUser, nil).Times(1)

				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), mockToken).Return(nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Dependencies
			mockRepo := userrepository.NewMockUserRepository(ctrl)
			mockTx := &mockTxManager{}
			mockToken, err := jwt.InitJWT(config.JWTConfig{
				SecretKey:  "mock_secret_key",
				RefreshKey: "mock_refresh_key",
			})
			if err != nil {
				t.Fatalf("InitJWT failed: %v", err)
			}

			// New User Usecase
			uc, err := userusecase.NewUserUsecase(&userusecase.UserUsecaseConfig{
				Repo:  mockRepo,
				Token: mockToken,
				Tx:    mockTx,
				DB:    nil,
			})
			if err != nil {
				t.Fatalf("NewUserUsecase failed: %v", err)
			}

			tc.mockFn(mockRepo)

			result, err := uc.RefreshToken(context.Background(), tc.input)

			// Assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestLogoutUsecase(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		mockFn      func(mockRepo *userrepository.MockUserRepository)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success revoked token",
			input: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), "mock_refresh_token").Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail token not found",
			input: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), "mock_refresh_token").Return(errs.ErrTokenNotFound).Times(1)
			},
			expectedErr: errs.ErrTokenNotFound,
		},
		{
			name:  "fail revoked token",
			input: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), "mock_refresh_token").Return(errors.New("db error")).Times(1)
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Dependencies
			mockRepo := userrepository.NewMockUserRepository(ctrl)
			mockTx := &mockTxManager{}
			mockToken, err := jwt.InitJWT(config.JWTConfig{
				SecretKey:  "mock_secret_key",
				RefreshKey: "mock_refresh_key",
			})
			if err != nil {
				t.Fatalf("InitJWT failed: %v", err)
			}

			// New User Usecase
			uc, err := userusecase.NewUserUsecase(&userusecase.UserUsecaseConfig{
				Repo:  mockRepo,
				Token: mockToken,
				Tx:    mockTx,
				DB:    nil,
			})
			if err != nil {
				t.Fatalf("NewUserUsecase failed: %v", err)
			}

			tc.mockFn(mockRepo)

			err = uc.Logout(context.Background(), tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
