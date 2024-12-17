package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/server/middleware"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestParseJsonBody(t *testing.T) {
	testcases := []struct {
		title      string
		body       string
		want       map[string]interface{}
		httpStatus int
	}{
		{
			title:      "Parse Json Body Success",
			body:       `{"key": "value"}`,
			want:       map[string]interface{}{"key": "value"},
			httpStatus: http.StatusOK,
		},
		{
			title:      "Parse Json Body Nested JSON",
			body:       `{"key": {"nestedKey": "nestedValue"}}`,
			want:       map[string]interface{}{"key": map[string]interface{}{"nestedKey": "nestedValue"}},
			httpStatus: http.StatusOK,
		},
		{
			title:      "Parse Json Body Empty Body",
			body:       "",
			want:       nil,
			httpStatus: http.StatusBadRequest,
		},
		{
			title:      "Parse Json Body Empty JSON",
			body:       "{}",
			want:       map[string]interface{}{},
			httpStatus: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/dummy", bytes.NewBufferString(tc.body))

			next := func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
				parsedPayload := req.Context().Value(middleware.PayloadContextKey).(map[string]interface{})
				assert.Equal(t, tc.want, parsedPayload)
			}

			handler := middleware.ParseJsonBody(next)
			handler(res, req, nil)

			assert.Equal(t, tc.httpStatus, res.Code)
		})
	}
}
