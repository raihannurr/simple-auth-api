package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/server/middleware"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestParseSession(t *testing.T) {
	cfg := config.AppConfig{
		Session: config.SessionConfig{
			SecretKey: "dummy-secret-key",
		},
	}
	expiredAt := time.Now().Add(10 * time.Minute).Unix()
	validCookie := utils.EncryptSessionValue(utils.Session{
		CSRFToken: utils.CSRFToken{
			Token:     "dummy-token",
			ExpiredAt: expiredAt,
		},
	}, cfg.Session)

	testcases := []struct {
		title                string
		sessionCookieValue   string
		expectedSessionValue utils.Session
	}{
		{
			title:              "Parse Session Success",
			sessionCookieValue: validCookie,
			expectedSessionValue: utils.Session{
				CSRFToken: utils.CSRFToken{
					Token:     "dummy-token",
					ExpiredAt: expiredAt,
				},
			},
		},
		{
			title:                "Empty Session Not Provided",
			sessionCookieValue:   "",
			expectedSessionValue: utils.Session{},
		},
		{
			title:                "Empty Session Invalid Value",
			sessionCookieValue:   "invalid-value",
			expectedSessionValue: utils.Session{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/dummy", nil)
			req.AddCookie(&http.Cookie{
				Name:  "session",
				Value: tc.sessionCookieValue,
			})

			next := func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
				session := req.Context().Value(middleware.SessionContextKey)
				assert.Equal(t, tc.expectedSessionValue, session)
			}

			sessionParser := middleware.NewSessionParser(cfg)
			handler := sessionParser(next)

			assert.NotPanics(t, func() {
				handler(res, req, nil)
			})
		})
	}
}
