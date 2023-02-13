package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
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

	log.Println("Starting server on port ", b.config.Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Println("Error on starting Server")
	} else {
		log.Fatal("Server Stopped F !")
	}
}
