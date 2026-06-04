package services

import (
	"fmt"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/repository/buddy_package_repo"
	"github.com/bored-bots/bored-template/internal/repository/user_repo"
	"github.com/bored-bots/bored-template/internal/services/buddy_package_service"
	"github.com/bored-bots/bored-template/internal/services/email_service"
	"github.com/bored-bots/bored-template/internal/services/health_service"
	"github.com/bored-bots/bored-template/internal/services/stripe_service"
	"github.com/bored-bots/bored-template/internal/services/user_service"
	"github.com/bored-bots/bored-template/internal/store"
	"github.com/stripe/stripe-go/v85"
)

type Services struct {
	HealthService       *health_service.Service
	BuddyPackageService *buddy_package_service.Service
	StripeService       *stripe_service.Service
	UserService         *user_service.Service
	EmailService        *email_service.Service
}

func NewService(config *config.Config, store *store.Store, stripeClient *stripe.Client) (*Services, error) {
	buddyPackageRepo := buddy_package_repo.New(store)
	userRepo := user_repo.New(store)

	emailService, err := email_service.New(&config.SMTP)
	if err != nil {
		return nil, fmt.Errorf("failed to create email service: %w", err)
	}

	return &Services{
		HealthService:       health_service.New(store),
		BuddyPackageService: buddy_package_service.New(buddyPackageRepo),
		StripeService:       stripe_service.New(buddyPackageRepo, userRepo, stripeClient, &config.Stripe),
		UserService:         user_service.New(userRepo, buddyPackageRepo),
		EmailService:        emailService,
	}, nil
}
