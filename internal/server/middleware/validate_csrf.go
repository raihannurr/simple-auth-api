package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raihannurr/simple-auth-api/internal/utils"
)

func ValidateCSRF(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		session := ctx.Value(SessionContextKey).(utils.Session)
		payload := ctx.Value(PayloadContextKey).(map[string]interface{})

		csrfToken, _ := payload["csrf_token"].(string)

		if session.CSRFToken.IsExpired() || session.CSRFToken.Token != csrfToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid CSRF token"))
			return
		}

		next(w, r, ps)
	}
}
