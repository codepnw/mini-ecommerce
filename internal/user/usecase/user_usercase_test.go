package userusecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/codepnw/mini-ecommerce/internal/user"
	userrepository "github.com/codepnw/mini-ecommerce/internal/user/repository"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
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
		mockFn      func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success",
			input: &user.User{Email: "user@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User) {
				u := mockUserData()
				mockRepo.EXPECT().Insert(gomock.Any(), nil, input).Return(u, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail email already exists",
			input: &user.User{Email: "user@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User) {
				mockRepo.EXPECT().Insert(gomock.Any(), nil, input).Return(nil, errs.ErrEmailAlreadyExists).Times(1)
			},
			expectedErr: errs.ErrEmailAlreadyExists,
		},
		{
			name:  "fail save token",
			input: &user.User{Email: "user@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User) {
				u := mockUserData()
				mockRepo.EXPECT().Insert(gomock.Any(), nil, input).Return(u, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockRepo, mockTx := setup(t)

			tc.mockFn(mockRepo, mockTx, tc.input)

			// Register Usecase
			result, err := uc.Register(context.Background(), tc.input)

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
		mockFn      func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success",
			input: &user.User{Email: "user@example.com", Password: "correct_password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User) {
				hashedPassword, _ := password.HashedPassword(input.Password)
				u := mockUserData()
				u.Email = input.Email
				u.Password = hashedPassword

				mockRepo.EXPECT().FindByEmail(gomock.Any(), input.Email).Return(u, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail wrong password",
			input: &user.User{Email: "user@example.com", Password: "wrong_password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User) {
				u := mockUserData()
				mockRepo.EXPECT().FindByEmail(gomock.Any(), input.Email).Return(u, nil).Times(1)

				hashedPassword, _ := password.HashedPassword(u.Password)
				password.ComparePassword(hashedPassword, input.Password)
			},
			expectedErr: errs.ErrUserCredentials,
		},
		{
			name:  "fail email not found",
			input: &user.User{Email: "user2@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User) {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), input.Email).Return(nil, errs.ErrUserNotFound).Times(1)
			},
			expectedErr: errs.ErrUserCredentials,
		},
		{
			name:  "fail save token",
			input: &user.User{Email: "user@example.com", Password: "password"},
			mockFn: func(mockRepo *userrepository.MockUserRepository, mockTx *mockTxManager, input *user.User) {
				hashedPassword, _ := password.HashedPassword(input.Password)
				u := mockUserData()
				u.Password = hashedPassword

				mockRepo.EXPECT().FindByEmail(gomock.Any(), input.Email).Return(u, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockRepo, mockTx := setup(t)

			tc.mockFn(mockRepo, mockTx, tc.input)

			// Login Usecase
			result, err := uc.Login(context.Background(), tc.input)

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
		token       string
		mockFn      func(mockRepo *userrepository.MockUserRepository, token string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				u := mockUserData()
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), token).Return(u.ID, nil).Times(1)

				mockRepo.EXPECT().FindByID(gomock.Any(), u.ID).Return(u, nil).Times(1)

				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail user not found",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				u := mockUserData()
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), token).Return(u.ID, nil).Times(1)

				mockRepo.EXPECT().FindByID(gomock.Any(), u.ID).Return(nil, errs.ErrUserNotFound).Times(1)
			},
			expectedErr: errs.ErrUserNotFound,
		},
		{
			name:  "fail save token",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				u := mockUserData()
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), token).Return(u.ID, nil).Times(1)

				mockRepo.EXPECT().FindByID(gomock.Any(), u.ID).Return(u, nil).Times(1)

				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), nil, gomock.Any()).Return(errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockRepo, _ := setup(t)

			tc.mockFn(mockRepo, tc.token)

			// RefreshToken Usecase
			result, err := uc.RefreshToken(context.Background(), tc.token)

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
		token       string
		mockFn      func(mockRepo *userrepository.MockUserRepository, token string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail token not found",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(errs.ErrTokenNotFound).Times(1)
			},
			expectedErr: errs.ErrTokenNotFound,
		},
		{
			name:  "fail revoked token",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(errDBMock).Times(1)
			},
			expectedErr: errDBMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockRepo, _ := setup(t)

			tc.mockFn(mockRepo, tc.token)

			// Logout Usecase
			err := uc.Logout(context.Background(), tc.token)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		mockFn      func(mockRepo *userrepository.MockUserRepository, userID int64)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:   "success",
			userID: 10,
			mockFn: func(mockRepo *userrepository.MockUserRepository, userID int64) {
				u := mockUserData()
				mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(u, nil).Times(1)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockRepo, _ := setup(t)

			tc.mockFn(mockRepo, tc.userID)

			// GetUser Usecase
			result, err := uc.GetUser(context.Background(), tc.userID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr) || err.Error() == tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// ================= Helper ======================
// -----------------------------------------------
func setup(t *testing.T) (userusecase.UserUsecase, *userrepository.MockUserRepository, *mockTxManager) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userrepository.NewMockUserRepository(ctrl)
	mockTx := &mockTxManager{}
	mockToken, err := jwt.InitJWT(config.JWTConfig{
		SecretKey:  "mock_secret_key",
		RefreshKey: "mock_refresh_key",
	})
	if err != nil {
		t.Fatalf("init jwt failed: %v", err)
	}

	uc, err := userusecase.NewUserUsecase(&userusecase.UserUsecaseConfig{
		Repo:  mockRepo,
		Token: mockToken,
		Tx:    mockTx,
	})
	if err != nil {
		t.Fatalf("user usecase failed: %v", err)
	}

	return uc, mockRepo, mockTx
}

func mockUserData() *user.User {
	return &user.User{
		ID:       10,
		Email:    "example@mail.com",
		Password: "example_password",
		Role:     "user",
	}
}

var errDBMock = errors.New("database error")
