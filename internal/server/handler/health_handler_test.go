package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/server/handler"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	handler := handler.NewHealthHandler(config.AppConfig{Port: "8080"})

	request, _ := http.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()

	handler.Health(response, request, nil)

	assert.Equal(t, response.Code, http.StatusOK)
}
