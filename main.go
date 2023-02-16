package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/raulcv/goapiws/handlers"
	"github.com/raulcv/goapiws/middleware"
	"github.com/raulcv/goapiws/server"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loanginf .env file variables")
	}
	SERVER_PORT := os.Getenv("SERVER_PORT")
	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")
	DATABASE_URL := os.Getenv("DATABASE_URL_SUPA")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        SERVER_PORT,
		JWTSecret:   JWT_SECRET_KEY,
		DatabaseUrl: DATABASE_URL,
	})
	if err != nil {
		log.Fatalf("Error creating server %v\n", err)
	}
	s.StartServer(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {

	r.Use(middleware.CheckAuthMiddleware(s))

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
}
