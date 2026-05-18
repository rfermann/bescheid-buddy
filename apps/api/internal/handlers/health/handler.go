package health_handler

import (
	"net/http"

	"github.com/bored-bots/bored-template/internal/services/health"
	"github.com/danielgtaylor/huma/v2"
)

type HealthHandler struct {
	healthService *health.Service
}

func NewHealthHandler(healthService *health.Service) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h *HealthHandler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check endpoint",
		Description: "Returns the health status of the API",
		Tags:        []string{"Health"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.healthCheck)
}
