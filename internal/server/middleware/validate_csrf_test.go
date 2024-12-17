package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/server/middleware"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestValidateCSRF(t *testing.T) {
	testcases := []struct {
		title          string
		sessionContext any
		payloadContext any
		httpStatus     int
	}{
		{
			title: "Validate CSRF Success",
			sessionContext: utils.Session{
				CSRFToken: utils.CSRFToken{
					Token:     "dummy-token",
					ExpiredAt: time.Now().Add(10 * time.Minute).Unix(),
				},
			},
			payloadContext: map[string]interface{}{"csrf_token": "dummy-token"},
			httpStatus:     http.StatusOK,
		},
		{
			title: "Validate CSRF Failed Expired",
			sessionContext: utils.Session{
				CSRFToken: utils.CSRFToken{
					Token:     "dummy-token",
					ExpiredAt: time.Now().Add(-10 * time.Minute).Unix(),
				},
			},
			payloadContext: map[string]interface{}{"csrf_token": "dummy-token"},
			httpStatus:     http.StatusUnauthorized,
		},
		{
			title: "Validate CSRF Failed Mismatched Token",
			sessionContext: utils.Session{
				CSRFToken: utils.CSRFToken{
					Token:     "dummy-token",
					ExpiredAt: time.Now().Add(10 * time.Minute).Unix(),
				},
			},
			payloadContext: map[string]interface{}{"csrf_token": "not-matched-token"},
			httpStatus:     http.StatusUnauthorized,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/dummy", nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.SessionContextKey, tc.sessionContext))
			req = req.WithContext(context.WithValue(req.Context(), middleware.PayloadContextKey, tc.payloadContext))

			next := func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
				res.WriteHeader(http.StatusOK)
			}

			handler := middleware.ValidateCSRF(next)
			handler(res, req, nil)

			assert.Equal(t, tc.httpStatus, res.Code)
		})
	}
}
