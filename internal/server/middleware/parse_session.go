package middleware

import (
	"context"
	"net/http"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/julienschmidt/httprouter"
)

func ParseSession(cfg config.SessionConfig, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session := utils.GetSessionCookie(r, cfg)
		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}

func NewSessionParser(cfg config.AppConfig) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return ParseSession(cfg.Session, next)
	}
}
