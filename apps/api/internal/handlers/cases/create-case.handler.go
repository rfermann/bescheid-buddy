package cases_handler

import (
	"context"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreateCaseInputBody struct {
	Title string `json:"title" maxLength:"255"`
}

type CreateCaseInput struct {
	Body CreateCaseInputBody `json:"body"`
}

type CreateCaseOutputBody struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type CreateCaseOutput struct {
	Body CreateCaseOutputBody `json:"body"`
}

func (h *Handler) createCase(ctx context.Context, input *CreateCaseInput) (*CreateCaseOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	newCase, err := h.service.CreateCase(ctx, cases_service.CreateCaseInput{
		UserID: userID,
		Title:  input.Body.Title,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create case", err)
	}

	return &CreateCaseOutput{
		Body: CreateCaseOutputBody{
			ID:    newCase.ID,
			Title: newCase.Title.String,
		},
	}, nil
}
