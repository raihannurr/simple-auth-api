package handler

import (
	"encoding/json"
	"net/http"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/errors"
	"github.com/raihannurr/simple-auth-api/internal/repository"
	"github.com/raihannurr/simple-auth-api/internal/server/middleware"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/julienschmidt/httprouter"
)

type AuthHandler struct {
	Config config.AppConfig
	Repo   repository.IRepository
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	payload := r.Context().Value(middleware.PayloadContextKey).(map[string]interface{})
	username, _ := payload["username"].(string)
	password, _ := payload["password"].(string)

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username and password are required"))
		return
	}

	user, err := h.Repo.GetUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(errors.ErrInvalidLoginCredentials.Error()))
		return
	}

	if !utils.VerifyPassword(password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(errors.ErrInvalidLoginCredentials.Error()))
		return
	}

	token := utils.GenerateToken(user, h.Config.JWT)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	payload := r.Context().Value(middleware.PayloadContextKey).(map[string]interface{})
	username, _ := payload["username"].(string)
	email, _ := payload["email"].(string)
	password, _ := payload["password"].(string)

	if username == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username, email, and password are required"))
		return
	}

	if !utils.IsStrongPassword(password) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.ErrInvalidPassword.Error()))
		return
	}

	user, err := h.Repo.CreateUser(username, email, password)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func NewAuthHandler(cfg config.AppConfig, repo repository.IRepository) AuthHandler {
	return AuthHandler{Config: cfg, Repo: repo}
}
