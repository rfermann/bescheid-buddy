package models_handler

import (
	"context"
	"errors"

	"github.com/bored-bots/bored-template/internal/services/models_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
)

type CreateModelInputBody struct {
	ProviderModelID string `json:"provider_model_id" required:"true" minLength:"1" maxLength:"255"`
	Name            string `json:"name" required:"true" minLength:"1" maxLength:"255"`
	Provider        string `json:"provider" required:"true" minLength:"1" maxLength:"255"`
	SortOrder       int32  `json:"sort_order"`
}

type CreateModelInput struct {
	Body CreateModelInputBody `json:"body"`
}

type CreateModelOutputBody struct {
	Model *db.AiModel `json:"model"`
}

type CreateModelOutput struct {
	Body CreateModelOutputBody `json:"body"`
}

func (h *Handler) createModel(ctx context.Context, input *CreateModelInput) (*CreateModelOutput, error) {
	model, err := h.service.CreateModel(ctx, models_service.CreateModelInput{
		ProviderModelID: input.Body.ProviderModelID,
		Name:            input.Body.Name,
		Provider:        input.Body.Provider,
		SortOrder:       input.Body.SortOrder,
	})
	if err != nil {
		switch {
		case errors.Is(err, models_service.ErrInvalidModel), errors.Is(err, models_service.ErrUnsupportedProvider):
			return nil, huma.Error422UnprocessableEntity(err.Error())
		case errors.Is(err, models_service.ErrModelAlreadyExists):
			return nil, huma.Error409Conflict(err.Error())
		default:
			return nil, huma.Error500InternalServerError("Failed to create AI model", err)
		}
	}

	return &CreateModelOutput{
		Body: CreateModelOutputBody{Model: model},
	}, nil
}
