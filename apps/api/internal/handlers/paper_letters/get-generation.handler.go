package paper_letters_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetGenerationInput struct {
	GenerationID uuid.UUID `path:"generationId"`
}

type GetGenerationOutput struct {
	Body GenerationBody
}

func (h *Handler) getGeneration(ctx context.Context, input *GetGenerationInput) (*GetGenerationOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	result, err := h.service.GetGeneration(ctx, userID, input.GenerationID)
	if err != nil {
		if errors.Is(err, paper_letters_service.ErrGenerationNotFound) {
			return nil, huma.Error404NotFound("Paper letter generation not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get generation", err)
	}
	body, err := generationBody(result)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to encode generation", err)
	}
	return &GetGenerationOutput{Body: body}, nil
}
