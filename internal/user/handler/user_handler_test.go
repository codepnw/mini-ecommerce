package userhandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codepnw/mini-ecommerce/internal/errs"
	userhandler "github.com/codepnw/mini-ecommerce/internal/user/handler"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		input          any
		mockFn         func(mockUc *userusecase.MockUserUsecase)
		expectedStatus int
	}

	testCases := []testCase{
		{
			name:  "success 200 OK",
			input: &userhandler.UserLoginReq{Email: "test@mail.com", Password: "test_password"},
			mockFn: func(mockUc *userusecase.MockUserUsecase) {
				mockUc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "fail 400 user credential",
			input: &userhandler.UserLoginReq{Email: "test@mail.com", Password: "123456"},
			mockFn: func(mockUc *userusecase.MockUserUsecase) {
				mockUc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errs.ErrUserCredentials).Times(1)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "fail 500 sever error",
			input: &userhandler.UserLoginReq{Email: "test@mail.com", Password: "123123"},
			mockFn: func(mockUc *userusecase.MockUserUsecase) {
				mockUc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errors.New("database timeout")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := userusecase.NewMockUserUsecase(ctrl)
			handler := userhandler.NewUserHandler(mockUC)
			router := gin.New()
			router.POST("/login", handler.Login)
			rr := httptest.NewRecorder()

			tc.mockFn(mockUC)

			body, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
