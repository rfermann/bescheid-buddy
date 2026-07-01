package services

import (
	"fmt"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/repository/buddy_packages_repository"
	"github.com/bored-bots/bored-template/internal/repository/support_cases_repository"
	"github.com/bored-bots/bored-template/internal/repository/users_repository"
	"github.com/bored-bots/bored-template/internal/services/buddy_packages_service"
	"github.com/bored-bots/bored-template/internal/services/email_service"
	"github.com/bored-bots/bored-template/internal/services/health_service"
	"github.com/bored-bots/bored-template/internal/services/stripe_service"
	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/bored-bots/bored-template/internal/services/users_service"
	"github.com/bored-bots/bored-template/internal/store"
	"github.com/stripe/stripe-go/v85"
)

type Services struct {
	HealthService       *health_service.Service
	BuddyPackageService *buddy_packages_service.Service
	StripeService       *stripe_service.Service
	UserService         *users_service.Service
	EmailService        *email_service.Service
	SupportCasesService *support_cases_service.Service
}

func NewService(config *config.Config, store *store.Store, stripeClient *stripe.Client) (*Services, error) {
	buddyPackageRepo := buddy_packages_repository.New(store)
	supportCasesRepo := support_cases_repository.New(store)
	userRepo := users_repository.New(store)

	emailService, err := email_service.New(&config.SMTP)
	if err != nil {
		return nil, fmt.Errorf("failed to create email service: %w", err)
	}

	return &Services{
		HealthService:       health_service.New(store),
		BuddyPackageService: buddy_packages_service.New(buddyPackageRepo),
		StripeService:       stripe_service.New(buddyPackageRepo, userRepo, stripeClient, &config.Stripe),
		UserService:         users_service.New(userRepo, buddyPackageRepo),
		EmailService:        emailService,
		SupportCasesService: support_cases_service.New(supportCasesRepo),
	}, nil
}
