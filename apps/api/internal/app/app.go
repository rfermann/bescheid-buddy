package app

import (
	"log"
	"os"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/server"
	"github.com/bored-bots/bored-template/internal/services"
	s "github.com/bored-bots/bored-template/internal/store"
	"github.com/stripe/stripe-go/v85"
)

type App struct {
	config *config.Config
	server *server.Server
}

func New() (*App, error) {
	cfg := loadConfig()

	store, err := s.NewStore(cfg.Database)
	if err != nil {
		return nil, err
	}

	stripeClient := stripe.NewClient(cfg.Stripe.SecretKey)
	svc, err := services.NewService(cfg, store, stripeClient)
	if err != nil {
		return nil, err
	}
	srv := server.NewServer(cfg, svc, store)

	return &App{
		config: cfg,
		server: srv,
	}, nil
}

func (a *App) Run() error {
	return a.server.Start()
}

func loadConfig() *config.Config {
	var configPath string
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		configPath = env
	} else {
		log.Fatalf("CONFIG_PATH environment variable is not set")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return cfg
}
