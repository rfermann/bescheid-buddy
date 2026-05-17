package services

import "github.com/bored-bots/bored-template/internal/store"

type Services struct {
	HealthService *HealthService
}

func NewService(store *store.Store) *Services {
	return &Services{
		HealthService: NewHealthService(store),
	}
}
