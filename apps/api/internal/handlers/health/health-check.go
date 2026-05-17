package health_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type healthCheckResponseBody struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Database  string `json:"database"`
}

type HealthCheckResponse struct {
	Body healthCheckResponseBody `json:"body"`
}

func (h *HealthHandler) healthCheck(ctx context.Context, input *struct{}) (*HealthCheckResponse, error) {
	result, err := h.healthService.Check(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Health check failed", err)
	}

	healthCheckResponse := healthCheckResponseBody{
		Status:    result.Status,
		Timestamp: result.Timestamp,
		Database:  result.Database,
	}

	return &HealthCheckResponse{
		Body: healthCheckResponse,
	}, nil
}
