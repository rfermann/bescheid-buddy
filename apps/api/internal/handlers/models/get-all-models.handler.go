package models_handler

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
)

type GetAllModelsOutputBody struct {
	Models []*db.AiModel `json:"models"`
}

type GetAllModelsOutput struct {
	Body GetAllModelsOutputBody `json:"body"`
}

func (h *Handler) getAllModels(ctx context.Context, _ *struct{}) (*GetAllModelsOutput, error) {
	models, err := h.service.GetModels(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get AI models", err)
	}

	return &GetAllModelsOutput{
		Body: GetAllModelsOutputBody{Models: models},
	}, nil
}
