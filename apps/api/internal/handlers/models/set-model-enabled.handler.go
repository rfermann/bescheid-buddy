package models_handler

import (
	"context"
	"errors"

	"github.com/bored-bots/bored-template/internal/services/models_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type SetModelEnabledInputBody struct {
	Enabled bool `json:"enabled" required:"true"`
}

type SetModelEnabledInput struct {
	ID   uuid.UUID                `path:"id"`
	Body SetModelEnabledInputBody `json:"body"`
}

type SetModelEnabledOutputBody struct {
	Model *db.AiModel `json:"model"`
}

type SetModelEnabledOutput struct {
	Body SetModelEnabledOutputBody `json:"body"`
}

func (h *Handler) setModelEnabled(ctx context.Context, input *SetModelEnabledInput) (*SetModelEnabledOutput, error) {
	model, err := h.service.SetModelEnabled(ctx, input.ID, input.Body.Enabled)
	if err != nil {
		if errors.Is(err, models_service.ErrModelNotFound) {
			return nil, huma.Error404NotFound(err.Error())
		}
		return nil, huma.Error500InternalServerError("Failed to set AI model enabled state", err)
	}

	return &SetModelEnabledOutput{
		Body: SetModelEnabledOutputBody{Model: model},
	}, nil
}
