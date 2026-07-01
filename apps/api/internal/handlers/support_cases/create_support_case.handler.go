package support_cases_handler

import (
	"context"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
)

type CreateSupportCaseInputBody struct {
	Subject  string                 `json:"subject" required:"true" maxLength:"255"`
	Category db.SupportCaseCategory `json:"category" required:"true" enum:"general,technical,billing,feature_request"`
	Message  string                 `json:"message" required:"true" maxLength:"255"`
}

type CreateSupportCaseInput struct {
	Body CreateSupportCaseInputBody `json:"body"`
}

type CreateSupportCaseOutputBody struct {
	Case *db.SupportCase `json:"case"`
}

type CreateSupportCaseOutput struct {
	Body CreateSupportCaseOutputBody `json:"body"`
}

func (h *Handler) createSupportCase(ctx context.Context, input *CreateSupportCaseInput) (*CreateSupportCaseOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	params := support_cases_service.CreateSupportCaseInput{
		UserID:   userID,
		Subject:  input.Body.Subject,
		Message:  input.Body.Message,
		Category: input.Body.Category,
	}

	supportCase, err := h.service.CreateSupportCase(ctx, params)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create support case", err)
	}

	return &CreateSupportCaseOutput{
		Body: CreateSupportCaseOutputBody{
			Case: supportCase,
		},
	}, nil
}
