package health_handler

import (
	"net/http"

	"github.com/bored-bots/bored-template/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

type HealthHandler struct {
	healthService *services.HealthService
}

func NewHealthHandler(healthService *services.HealthService) *HealthHandler {
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
	}, h.healthCheck)
}
