package models_handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-models",
		Method:      http.MethodGet,
		Path:        "/models",
		Summary:     "Get selectable AI models",
		Tags:        []string{"Models"},
	}, h.getModels)
}
