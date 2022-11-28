package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"githuh.com/go/rest-crud/database"
	"githuh.com/go/rest-crud/repository"
	websocket "githuh.com/go/rest-crud/websocket"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (this *Broker) Config() *Config {
	return this.config
}

func (this *Broker) Hub() *websocket.Hub {
	return this.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("Secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("Database url is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

func (this *Broker) Start(binder func(s Server, r *mux.Router)) {
	this.router = mux.NewRouter()
	binder(this, this.router)
	handler := cors.Default().Handler(this.router)
	repo, err := database.NewPostgresRespository(this.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go this.hub.Run()
	repository.SetRepository(repo)
	log.Println("Serve on port", this.Config().Port)
	if err := http.ListenAndServe(this.Config().Port, handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
