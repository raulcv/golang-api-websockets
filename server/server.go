package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/raulcv/goapiws/database"
	"github.com/raulcv/goapiws/repository"
	"github.com/rs/cors"
)

type Config struct {
	Host        string
	Port        string
	JWTSecret   string
	DatabaseUrl string
}
type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port Number is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("Secret Key is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("Database url is required")
	}
	Broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return Broker, nil
}

func (b *Broker) StartServer(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	handler := cors.Default().Handler(b.router)

	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	repository.SetRepository(repo)

	srv := &http.Server{
		Handler: handler,
		Addr:    b.config.Host + b.config.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server on port ", b.config.Port, " | open in: http://"+srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Error on starting Server", err)
	} else {
		log.Fatal("Server Stopped F !")
	}
}
