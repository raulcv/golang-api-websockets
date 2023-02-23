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
		log.Fatal("Error Loading .env file variables")
	}
	SERVER_HOST := os.Getenv("SERVER_HOST")
	SERVER_PORT := os.Getenv("SERVER_PORT")
	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")
	DATABASE_URL := os.Getenv("DATABASE_URL_SUPA")

	s, err := server.NewServer(context.Background(), &server.Config{
		Host:        SERVER_HOST,
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

	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.CheckAuthMiddleware(s))

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts", handlers.AddPostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
	api.HandleFunc("/posts/{id}", handlers.ActivatePostHandler(s)).Methods(http.MethodPatch)
	r.HandleFunc("/posts", handlers.ListPostHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/websoscket", s.Hub().HandleWebSocket)
}
