package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/errors"
	mock_repository "github.com/raihannurr/simple-auth-api/internal/repository/mocks"
	"github.com/raihannurr/simple-auth-api/internal/server/handler"
	"github.com/raihannurr/simple-auth-api/internal/server/middleware"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		title               string
		userContext         any
		expectedCode        int
		mockRepoExpectation func(mockRepo *mock_repository.MockIRepository)
		expectedBody        string
		expectedError       error
	}{
		{
			title:        "Success",
			userContext:  uint(123),
			expectedCode: http.StatusOK,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByID(uint(123)).Return(
					entity.User{ID: 123, Username: "user1", Email: "user1@example.com", Verified: false, Description: ""},
					nil,
				)
			},
			expectedBody: `{"id":123,"username":"user1","email":"user1@example.com","verified":false,"description":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			title:        "Failed user not found",
			userContext:  uint(123),
			expectedCode: http.StatusNotFound,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByID(uint(123)).Return(
					entity.User{},
					errors.ErrUserNotFound,
				)
			},
			expectedError: errors.ErrUserNotFound,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.title, func(t *testing.T) {
			mockRepo := mock_repository.NewMockIRepository(ctrl)
			testcase.mockRepoExpectation(mockRepo)

			ctx := context.WithValue(context.Background(), middleware.UserIDContextKey, testcase.userContext)
			request, _ := http.NewRequest("GET", "/profile", nil)
			request = request.WithContext(ctx)
			response := httptest.NewRecorder()

			handler := handler.NewUserHandler(config.AppConfig{}, mockRepo)
			handler.GetProfile(response, request, nil)

			assert.Equal(t, response.Code, testcase.expectedCode)
			if testcase.expectedError != nil {
				assert.Equal(t, response.Body.String(), testcase.expectedError.Error())
			} else {
				assert.JSONEq(t, response.Body.String(), testcase.expectedBody)
			}
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		title               string
		userContext         any
		payload             any
		expectedCode        int
		mockRepoExpectation func(mockRepo *mock_repository.MockIRepository)
		expectedBody        string
		expectedError       error
	}{
		{
			title:        "Success",
			userContext:  uint(123),
			payload:      map[string]interface{}{"description": "new description"},
			expectedCode: http.StatusAccepted,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByID(uint(123)).Return(
					entity.User{ID: 123, Username: "user1", Email: "user1@example.com", Verified: false, Description: ""},
					nil,
				)
				mockRepo.EXPECT().UpdateUser(uint(123), entity.User{ID: 123, Username: "user1", Email: "user1@example.com", Verified: false, Description: "new description"}).Return(
					entity.User{ID: 123, Username: "user1", Email: "user1@example.com", Verified: false, Description: "new description"},
					nil,
				)
			},
			expectedBody: `{"id":123,"username":"user1","email":"user1@example.com","verified":false,"description":"new description","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			title:        "Failed user not found",
			userContext:  uint(123),
			payload:      map[string]interface{}{"description": "new description"},
			expectedCode: http.StatusNotFound,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByID(uint(123)).Return(
					entity.User{},
					errors.ErrUserNotFound,
				)
			},
			expectedError: errors.ErrUserNotFound,
		},
		{
			title:        "Failed internal server error",
			userContext:  uint(123),
			payload:      map[string]interface{}{"description": "new description"},
			expectedCode: http.StatusInternalServerError,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByID(uint(123)).Return(
					entity.User{ID: 123, Username: "user1", Email: "user1@example.com", Verified: false, Description: ""},
					nil,
				)
				mockRepo.EXPECT().UpdateUser(uint(123), entity.User{ID: 123, Username: "user1", Email: "user1@example.com", Verified: false, Description: "new description"}).Return(
					entity.User{},
					errors.ErrInternalServerError,
				)
			},
			expectedError: errors.ErrInternalServerError,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.title, func(t *testing.T) {
			mockRepo := mock_repository.NewMockIRepository(ctrl)
			testcase.mockRepoExpectation(mockRepo)

			ctx := context.WithValue(context.Background(), middleware.UserIDContextKey, testcase.userContext)
			ctx = context.WithValue(ctx, middleware.PayloadContextKey, testcase.payload)
			request, _ := http.NewRequest("PATCH", "/profile", nil)
			request = request.WithContext(ctx)
			response := httptest.NewRecorder()

			handler := handler.NewUserHandler(config.AppConfig{}, mockRepo)
			handler.UpdateProfile(response, request, nil)

			assert.Equal(t, response.Code, testcase.expectedCode)
			if testcase.expectedError != nil {
				assert.Equal(t, response.Body.String(), testcase.expectedError.Error())
			} else {
				assert.JSONEq(t, response.Body.String(), testcase.expectedBody)
			}
		})
	}
}
