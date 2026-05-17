package services

import (
	"context"
	"time"

	"github.com/bored-bots/bored-template/internal/store"
)

type HealthService struct {
	store *store.Store
}

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Database  string `json:"database"`
}

func NewHealthService(store *store.Store) *HealthService {
	return &HealthService{
		store: store,
	}
}

func (s *HealthService) Check(ctx context.Context) (*HealthCheckResponse, error) {
	status := &HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Database:  "healthy",
	}

	if err := s.checkDatabase(ctx); err != nil {
		status.Status = "unhealthy"
		status.Database = "unhealthy"
		return status, nil
	}

	return status, nil
}

func (s *HealthService) checkDatabase(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	return s.store.Ping(ctx)
}
