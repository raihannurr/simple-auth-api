package main

import (
	"log"
	"net/http"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/server"
)

func main() {
	log.Println("Starting server ...")
	cfg := config.LoadEnv(".env")

	if cfg.Port == "" {
		log.Fatal("environment variable APP_PORT is not set")
	}

	router := server.NewRouter(cfg)

	log.Println("Server started on port:", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, router)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
