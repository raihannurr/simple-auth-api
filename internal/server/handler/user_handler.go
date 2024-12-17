package handler

import (
	"encoding/json"
	"net/http"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/repository"
	"github.com/raihannurr/simple-auth-api/internal/server/middleware"

	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	Config config.AppConfig
	Repo   repository.IRepository
}

func (h UserHandler) GetProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userContext := r.Context().Value(middleware.UserIDContextKey).(uint)
	user, err := h.Repo.GetUserByID(userContext)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (h UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userContext := r.Context().Value(middleware.UserIDContextKey).(uint)
	payload := r.Context().Value(middleware.PayloadContextKey).(map[string]interface{})
	user, err := h.Repo.GetUserByID(userContext)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	user.Description = payload["description"].(string)
	user, err = h.Repo.UpdateUser(userContext, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(user)
}

func NewUserHandler(cfg config.AppConfig, repo repository.IRepository) UserHandler {
	return UserHandler{Config: cfg, Repo: repo}
}
