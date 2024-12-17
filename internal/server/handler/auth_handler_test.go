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
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		title                string
		payload              any
		expectedCode         int
		mockRepoExpectation  func(mockRepo *mock_repository.MockIRepository)
		expectedBodyContains string
		expectedError        error
	}{
		{
			title:        "Success",
			payload:      map[string]interface{}{"username": "user1", "password": "password"},
			expectedCode: http.StatusOK,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByUsername("user1").Return(
					entity.User{ID: 123, Password: utils.HashPassword("password")},
					nil,
				)
			},
			expectedBodyContains: `{"token":`,
		},
		{
			title:               "Failed empty request body",
			payload:             map[string]interface{}{"username": "", "password": ""},
			expectedCode:        http.StatusBadRequest,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {},
			expectedError:       errors.New("Username and password are required"),
		},
		{
			title:        "Failed user not found",
			payload:      map[string]interface{}{"username": "user1", "password": "password"},
			expectedCode: http.StatusUnauthorized,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByUsername("user1").Return(
					entity.User{},
					errors.ErrUserNotFound,
				)
			},
			expectedError: errors.ErrInvalidLoginCredentials,
		},
		{
			title:        "Failed password not match",
			payload:      map[string]interface{}{"username": "user1", "password": "password"},
			expectedCode: http.StatusUnauthorized,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().GetUserByUsername("user1").Return(
					entity.User{ID: 123, Password: utils.HashPassword("new-password")},
					nil,
				)
			},
			expectedError: errors.ErrInvalidLoginCredentials,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.title, func(t *testing.T) {
			mockRepo := mock_repository.NewMockIRepository(ctrl)
			testcase.mockRepoExpectation(mockRepo)

			ctx := context.WithValue(context.Background(), middleware.PayloadContextKey, testcase.payload)
			request, _ := http.NewRequest("POST", "/login", nil)
			request = request.WithContext(ctx)
			response := httptest.NewRecorder()

			handler := handler.NewAuthHandler(config.AppConfig{}, mockRepo)
			handler.Login(response, request, nil)

			assert.Equal(t, response.Code, testcase.expectedCode)
			if testcase.expectedError != nil {
				assert.Equal(t, response.Body.String(), testcase.expectedError.Error())
			} else {
				assert.Contains(t, response.Body.String(), testcase.expectedBodyContains)
			}
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		title               string
		payload             any
		expectedCode        int
		mockRepoExpectation func(mockRepo *mock_repository.MockIRepository)
		expectedBody        string
		expectedError       error
	}{
		{
			title:        "Success",
			payload:      map[string]interface{}{"username": "user1", "email": "user1@example.com", "password": "P4sssword!!!"},
			expectedCode: http.StatusCreated,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().CreateUser("user1", "user1@example.com", "P4sssword!!!").Return(
					entity.User{ID: 123, Username: "user1", Email: "user1@example.com", Description: "", Verified: false},
					nil,
				)
			},
			expectedBody: `{"id":123,"username":"user1","email":"user1@example.com","description":"","verified":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			title:               "Failed empty request body",
			payload:             map[string]interface{}{"username": "", "email": "", "password": ""},
			expectedCode:        http.StatusBadRequest,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {},
			expectedError:       errors.New("Username, email, and password are required"),
		},
		{
			title:               "Failed password not strong",
			payload:             map[string]interface{}{"username": "user1", "email": "user1@example.com", "password": "password"},
			expectedCode:        http.StatusBadRequest,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {},
			expectedError:       errors.ErrInvalidPassword,
		},
		{
			title:        "Failed create user error",
			payload:      map[string]interface{}{"username": "user1", "email": "user1@example.com", "password": "P4sssword!!!"},
			expectedCode: http.StatusUnprocessableEntity,
			mockRepoExpectation: func(mockRepo *mock_repository.MockIRepository) {
				mockRepo.EXPECT().CreateUser("user1", "user1@example.com", "P4sssword!!!").Return(
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

			ctx := context.WithValue(context.Background(), middleware.PayloadContextKey, testcase.payload)
			request, _ := http.NewRequest("POST", "/register", nil)
			request = request.WithContext(ctx)
			response := httptest.NewRecorder()

			handler := handler.NewAuthHandler(config.AppConfig{}, mockRepo)
			handler.Register(response, request, nil)

			assert.Equal(t, response.Code, testcase.expectedCode)
			if testcase.expectedError != nil {
				assert.Equal(t, response.Body.String(), testcase.expectedError.Error())
			} else {
				assert.JSONEq(t, response.Body.String(), testcase.expectedBody)
			}
		})
	}
}
