package models_handler

import (
	"context"

	"github.com/bored-bots/bored-template/internal/services/llm_service"
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

func (h *Handler) getModels(_ context.Context, _ *struct{}) (*GetModelsOutput, error) {
	models := llm_service.Models()
	result := make([]ModelBody, len(models))
	for i, model := range models {
		result[i] = ModelBody{
			ID:       string(model.ID),
			Name:     model.Name,
			Provider: model.Provider,
		}
	}

	return &GetModelsOutput{Body: GetModelsOutputBody{Models: result}}, nil
}
