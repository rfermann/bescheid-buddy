package services

import (
	"fmt"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/repository/buddy_packages_repository"
	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/repository/support_cases_repository"
	"github.com/bored-bots/bored-template/internal/repository/users_repository"
	"github.com/bored-bots/bored-template/internal/services/buddy_packages_service"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/bored-bots/bored-template/internal/services/document_service"
	"github.com/bored-bots/bored-template/internal/services/email_service"
	"github.com/bored-bots/bored-template/internal/services/health_service"
	"github.com/bored-bots/bored-template/internal/services/models_service"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_pdf_service"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/bored-bots/bored-template/internal/services/queue_service"
	"github.com/bored-bots/bored-template/internal/services/stripe_service"
	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/bored-bots/bored-template/internal/services/users_service"
	"github.com/bored-bots/bored-template/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/riverqueue/river"
	"github.com/stripe/stripe-go/v85"
)

type Services struct {
	HealthService          *health_service.Service
	BuddyPackageService    *buddy_packages_service.Service
	StripeService          *stripe_service.Service
	UserService            *users_service.Service
	EmailService           *email_service.Service
	SupportCasesService    *support_cases_service.Service
	CasesService           *cases_service.Service
	PaperLettersService    *paper_letters_service.Service
	PaperLettersPDFService *paper_letters_pdf_service.Service
	ModelsService          *models_service.Service
	DocumentService        *document_service.Service
	QueueService           *queue_service.Service
}

func NewService(config *config.Config, modelsService *models_service.Service, userService *users_service.Service, minioClient *minio.Client, riverClient *river.Client[pgx.Tx], store *store.Store, stripeClient *stripe.Client) (*Services, error) {
	buddyPackageRepo := buddy_packages_repository.New(store)
	supportCasesRepo := support_cases_repository.New(store)
	userRepo := users_repository.New(store)
	casesRepo := cases_repository.New(store)

	emailService, err := email_service.New(&config.SMTP)
	if err != nil {
		return nil, fmt.Errorf("failed to create email service: %w", err)
	}

	documentService := document_service.New(minioClient, config.S3)
	queueService := queue_service.New(riverClient)

	return &Services{
		HealthService:          health_service.New(store),
		BuddyPackageService:    buddy_packages_service.New(buddyPackageRepo),
		StripeService:          stripe_service.New(buddyPackageRepo, userRepo, stripeClient, &config.Stripe),
		UserService:            userService,
		EmailService:           emailService,
		SupportCasesService:    support_cases_service.New(supportCasesRepo),
		CasesService:           cases_service.New(casesRepo, documentService, queueService, modelsService, config.Upload.MaxBytes),
		PaperLettersService:    paper_letters_service.New(casesRepo, queueService, modelsService, userService, documentService, config.PaperLetter.MaxPDFBytes, config.Pingen),
		PaperLettersPDFService: paper_letters_pdf_service.New(),
		ModelsService:          modelsService,
		DocumentService:        documentService,
		QueueService:           queueService,
	}, nil
}
