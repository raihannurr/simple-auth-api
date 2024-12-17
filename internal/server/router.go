package server

import (
	"net/http"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/repository"
	"github.com/raihannurr/simple-auth-api/internal/server/handler"
	"github.com/raihannurr/simple-auth-api/internal/server/middleware"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(cfg config.AppConfig) http.Handler {
	repo := repository.InitializeRepository(cfg)
	healthHandler := handler.NewHealthHandler(cfg)
	csrfHandler := handler.NewCSRFHandler(cfg)
	authHandler := handler.NewAuthHandler(cfg, repo)
	userHandler := handler.NewUserHandler(cfg, repo)

	authenticate := middleware.NewAuthenticator(cfg)
	sessionParser := middleware.NewSessionParser(cfg)

	router := httprouter.New()

	router.GET("/", healthHandler.Health) // for health check
	router.GET("/csrf-token", csrfHandler.GetCSRFToken)

	// for Debugging
	router.POST("/csrf-verify",
		sessionParser(
			middleware.ParseJsonBody(
				middleware.ValidateCSRF(healthHandler.Health),
			),
		),
	)

	router.POST("/login",
		sessionParser(
			middleware.ParseJsonBody(authHandler.Login),
		),
	)

	router.POST("/register",
		sessionParser(
			middleware.ParseJsonBody(
				middleware.ValidateCSRF(authHandler.Register),
			),
		),
	)

	router.GET("/profile", authenticate(userHandler.GetProfile))
	router.PATCH("/profile", authenticate(
		middleware.ParseJsonBody(
			userHandler.UpdateProfile,
		),
	))

	return router
}
