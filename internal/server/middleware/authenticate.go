package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/julienschmidt/httprouter"
)

func Authenticate(cfg config.JWTConfig, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")
		tokenVal, err := utils.VerifyToken(token, cfg)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}

		userId, _ := strconv.Atoi(tokenVal.Subject)
		ctx := context.WithValue(r.Context(), UserIDContextKey, uint(userId))
		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}

func NewAuthenticator(cfg config.AppConfig) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return Authenticate(cfg.JWT, next)
	}
}
