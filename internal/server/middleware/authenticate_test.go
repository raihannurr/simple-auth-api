package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/server/middleware"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	cfg := config.AppConfig{
		JWT: config.JWTConfig{
			PrivateKey: "dummy-key",
			Lifetime:   1 * time.Hour,
		},
	}
	validToken1 := utils.GenerateToken(entity.User{ID: 123}, cfg.JWT)

	testcases := []struct {
		title          string
		authHeader     string
		expectedUserID uint
		httpStatus     int
	}{
		{
			title:          "Authenticate Success",
			authHeader:     validToken1,
			expectedUserID: 123,
			httpStatus:     http.StatusOK,
		},
		{
			title:      "Authenticate Failed Not Provided",
			httpStatus: http.StatusUnauthorized,
		},
		{
			title:      "Authenticate Failed Invalid Token",
			authHeader: "invalidTokenValue",
			httpStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/dummy", nil)
			req.Header.Set("Authorization", tc.authHeader)

			next := func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
				username := req.Context().Value(middleware.UserIDContextKey)
				assert.Equal(t, tc.expectedUserID, username)
			}

			authenticator := middleware.NewAuthenticator(cfg)
			handler := authenticator(next)
			handler(res, req, nil)

			assert.Equal(t, tc.httpStatus, res.Code)
		})
	}
}
