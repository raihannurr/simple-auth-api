package handler

import (
	"net/http"

	"github.com/raihannurr/simple-auth-api/internal/config"

	"github.com/julienschmidt/httprouter"
)

type HealthHandler struct {
	Config config.AppConfig
}

func (h HealthHandler) Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = w.Write([]byte("Simple Auth REST API at " + h.Config.Port))

}

func NewHealthHandler(config config.AppConfig) HealthHandler {
	return HealthHandler{Config: config}
}
