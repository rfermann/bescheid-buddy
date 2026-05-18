package services

import (
	"github.com/bored-bots/bored-template/internal/repository"
	"github.com/bored-bots/bored-template/internal/services/buddy_package"
	"github.com/bored-bots/bored-template/internal/services/health"
	"github.com/bored-bots/bored-template/internal/store"
)

type Services struct {
	HealthService       *health.Service
	BuddyPackageService *buddy_package.Service
}

func NewService(store *store.Store) *Services {
	buddyPackageRepo := repository.NewBuddyPackageRepository(store)

	return &Services{
		HealthService:       health.NewService(store),
		BuddyPackageService: buddy_package.NewService(buddyPackageRepo),
	}
}
