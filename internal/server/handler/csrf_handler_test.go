package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/server/handler"

	"github.com/stretchr/testify/assert"
)

func TestCSRFHandler_GetCSRFToken(t *testing.T) {
	handler := handler.NewCSRFHandler(config.AppConfig{
		Session: config.SessionConfig{
			SecretKey: "dummy-secret-key",
		},
	})

	request, _ := http.NewRequest("GET", "/csrf-token", nil)
	response := httptest.NewRecorder()

	handler.GetCSRFToken(response, request, nil)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Contains(t, response.Body.String(), "csrf_token")
}
