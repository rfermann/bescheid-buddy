package models_handler

import (
	"net/http"

	"github.com/bored-bots/bored-template/internal/middleware"
	"github.com/bored-bots/bored-template/internal/services/models_service"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *models_service.Service
}

func NewHandler(service *models_service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-models",
		Method:      http.MethodGet,
		Path:        "/models",
		Summary:     "Get selectable AI models",
		Tags:        []string{"Models"},
	}, h.getModels)

	huma.Register(api, huma.Operation{
		OperationID: "get-model-providers",
		Method:      http.MethodGet,
		Path:        "/admin/model-providers",
		Summary:     "Get supported AI model providers",
		Tags:        []string{"Models"},
		Security:    []map[string][]string{{"session-cookie": {middleware.RoleAdmin}}},
	}, h.getModelProviders)

	huma.Register(api, huma.Operation{
		OperationID: "get-all-models",
		Method:      http.MethodGet,
		Path:        "/admin/models",
		Summary:     "Get all AI models",
		Tags:        []string{"Models"},
		Security:    []map[string][]string{{"session-cookie": {middleware.RoleAdmin}}},
	}, h.getAllModels)

	huma.Register(api, huma.Operation{
		OperationID:   "create-ai-model",
		Method:        http.MethodPost,
		Path:          "/admin/models",
		Summary:       "Create an AI model",
		Tags:          []string{"Models"},
		Security:      []map[string][]string{{"session-cookie": {middleware.RoleAdmin}}},
		DefaultStatus: http.StatusCreated,
	}, h.createModel)

	huma.Register(api, huma.Operation{
		OperationID: "set-model-enabled",
		Method:      http.MethodPatch,
		Path:        "/admin/models/{id}",
		Summary:     "Activate or deactivate an AI model",
		Tags:        []string{"Models"},
		Security:    []map[string][]string{{"session-cookie": {middleware.RoleAdmin}}},
	}, h.setModelEnabled)
}
