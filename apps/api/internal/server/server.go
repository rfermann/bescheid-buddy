package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/handlers/buddy_packages"
	"github.com/bored-bots/bored-template/internal/handlers/cases"
	"github.com/bored-bots/bored-template/internal/handlers/email"
	"github.com/bored-bots/bored-template/internal/handlers/health"
	"github.com/bored-bots/bored-template/internal/handlers/models"
	"github.com/bored-bots/bored-template/internal/handlers/stripe"
	"github.com/bored-bots/bored-template/internal/handlers/support_cases"
	"github.com/bored-bots/bored-template/internal/handlers/users"
	"github.com/bored-bots/bored-template/internal/middleware"
	"github.com/bored-bots/bored-template/internal/services"
	"github.com/bored-bots/bored-template/internal/store"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type QueueStopper interface {
	Stop(ctx context.Context) error
}

type Server struct {
	httpServer *http.Server
	Router     chi.Router
	API        huma.API
	Config     *config.Config
	store      *store.Store
	queue      QueueStopper
}

func NewServer(cfg *config.Config, services *services.Services, store *store.Store, queue QueueStopper) *Server {
	router := chi.NewRouter()

	humaConfig := huma.DefaultConfig(cfg.API.Title, cfg.API.Version)
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"cookieAuth": {
			Type:        "http",
			Scheme:      "Cookie",
			Name:        "session-cookie",
			Description: "cookie based authentication using session tokens.",
		},
		"apiKeyAuth": {
			Type:        "apiKey",
			Name:        "X-API-Key",
			In:          "header",
			Description: "API key authentication using X-API-Key header.",
		},
	}
	api := humachi.New(router, humaConfig)
	api.UseMiddleware(middleware.NewAuthMiddleware(api, &cfg.Auth))
	apiV1 := huma.NewGroup(api, "/v1")

	healthHandler := health_handler.NewHealthHandler(services.HealthService)
	healthHandler.RegisterRoutes(apiV1)

	modelsHandler := models_handler.NewHandler()
	modelsHandler.RegisterRoutes(apiV1)

	buddyPackageHandler := buddy_packages_handler.NewHandler(services.BuddyPackageService)
	buddyPackageHandler.RegisterRoutes(api)

	stripeHandler := stripe_handler.NewHandler(services.StripeService)
	stripeHandler.RegisterRoutes(api)

	userHandler := users_handler.NewHandler(services.UserService)
	userHandler.RegisterRoutes(api)

	emailHandler := email_handler.NewHandler(services.EmailService)
	emailHandler.RegisterRoutes(api)

	supportCasesHandler := support_cases_handler.NewHandler(services.SupportCasesService)
	supportCasesHandler.RegisterRoutes(apiV1)

	casesHandler := cases_handler.NewHandler(services.CasesService)
	casesHandler.RegisterRoutes(apiV1)

	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.GetAddr(),
			Handler:      router,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		},
		Router: router,
		API:    api,
		Config: cfg,
		store:  store,
		queue:  queue,
	}
}

func (s *Server) Start() error {
	fmt.Printf("Starting server on %s\n", s.Config.GetAddr())

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	fmt.Println("Shutting down server…")

	fmt.Printf("GracefulShutdown: %d", s.Config.Server.GracefulShutdown)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.Config.Server.GracefulShutdown)*time.Second)
	defer cancel()
	defer s.store.Close()

	if s.queue != nil {
		if err := s.queue.Stop(ctx); err != nil {
			fmt.Printf("Queue failed to stop cleanly: %v\n", err)
		}
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
		return err
	}

	fmt.Println("Server exited")
	return nil
}
