package services

import (
	"github.com/bored-bots/bored-template/internal/services/health"
	"github.com/bored-bots/bored-template/internal/store"
)

type Services struct {
	HealthService *health.Service
}

func NewService(store *store.Store) *Services {
	return &Services{
		HealthService: health.NewService(store),
	}
}
