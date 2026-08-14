package app

import (
	"context"
	"log"
	"os"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/queue"
	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/server"
	"github.com/bored-bots/bored-template/internal/services"
	"github.com/bored-bots/bored-template/internal/services/analysis_service"
	"github.com/bored-bots/bored-template/internal/services/document_service"
	"github.com/bored-bots/bored-template/internal/services/extraction_service"
	"github.com/bored-bots/bored-template/internal/services/llm_service"
	"github.com/bored-bots/bored-template/internal/services/reply_service"
	"github.com/bored-bots/bored-template/internal/services/vision_service"
	s "github.com/bored-bots/bored-template/internal/store"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	minioClient, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		Region: cfg.S3.Region,
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	documentService := document_service.New(minioClient, cfg.S3)
	visionService := vision_service.New(cfg.Vision.APIKey)
	casesRepo := cases_repository.New(store)
	extractionService := extraction_service.New(documentService, visionService, casesRepo, cfg.Upload.MaxBytes)
	llmService := llm_service.New(llm_service.Config{
		OpenAIAPIKey:    cfg.LLM.OpenAIAPIKey,
		GeminiAPIKey:    cfg.LLM.GeminiAPIKey,
		AnthropicAPIKey: cfg.LLM.AnthropicAPIKey,
	})
	analysisService := analysis_service.New(casesRepo, llmService)
	replyService := reply_service.New(casesRepo, llmService)

	q, err := queue.New(store.Pool, extractionService, analysisService, replyService)
	if err != nil {
		return nil, err
	}
	if err := q.Start(context.Background()); err != nil {
		return nil, err
	}

	svc, err := services.NewService(cfg, minioClient, q.Client, store, stripeClient)
	if err != nil {
		return nil, err
	}
	srv := server.NewServer(cfg, svc, store, q)

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
