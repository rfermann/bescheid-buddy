package models_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type ModelBody struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type GetModelsOutputBody struct {
	Models []ModelBody `json:"models"`
}

type GetModelsOutput struct {
	Body GetModelsOutputBody `json:"body"`
}

func (h *Handler) getModels(ctx context.Context, _ *struct{}) (*GetModelsOutput, error) {
	models, err := h.service.GetEnabledModels(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get AI models", err)
	}

	result := make([]ModelBody, len(models))
	for i, model := range models {
		result[i] = ModelBody{
			ID:       model.ProviderModelID,
			Name:     model.Name,
			Provider: model.Provider,
		}
	}

	return &GetModelsOutput{Body: GetModelsOutputBody{Models: result}}, nil
}
