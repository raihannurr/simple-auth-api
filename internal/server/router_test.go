package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/server"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := server.NewRouter(config.AppConfig{}).(*httprouter.Router)
	assert.NotNil(t, router)

	testcases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/"},
		{http.MethodGet, "/csrf-token"},
		{http.MethodPost, "/login"},
		{http.MethodPost, "/register"},
		{http.MethodGet, "/profile"},
		{http.MethodPatch, "/profile"},
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("Test %s %s", testcase.method, testcase.path), func(t *testing.T) {
			handler, _, _ := router.Lookup(testcase.method, testcase.path)
			assert.NotNil(t, handler) // Validate that the route is handled
		})
	}
}
