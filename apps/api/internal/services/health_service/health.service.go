package health_service

import (
	"context"
	"time"

	"github.com/bored-bots/bored-template/internal/store"
)

type Service struct {
	store *store.Store
}

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Database  string `json:"database"`
}

func New(store *store.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Check(ctx context.Context) (*HealthCheckResponse, error) {
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

func (s *Service) checkDatabase(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	return s.store.Ping(ctx)
}
