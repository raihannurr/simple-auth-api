package handler

import (
	"fmt"
	"net/http"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/julienschmidt/httprouter"
)

type CSRFHandler struct {
	Config config.AppConfig
}

func (h CSRFHandler) GetCSRFToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	csrfToken := utils.GenerateCSRFToken(h.Config.Session)
	utils.SetSessionCookie(
		w,
		utils.Session{CSRFToken: csrfToken},
		h.Config.Session,
	)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"csrf_token": "%s"}`, csrfToken.Token)))
}

func NewCSRFHandler(config config.AppConfig) CSRFHandler {
	return CSRFHandler{Config: config}
}
