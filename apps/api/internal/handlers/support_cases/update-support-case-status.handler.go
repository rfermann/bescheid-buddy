package support_cases_handler

import (
	"context"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UpdateSupportCaseStatusInputBody struct {
	Status db.SupportCaseStatus `json:"status" required:"true" enum:"in_progress,resolved,closed"`
}

type UpdateSupportCaseStatusInput struct {
	ID   uuid.UUID                        `path:"id"`
	Body UpdateSupportCaseStatusInputBody `json:"body"`
}

type UpdateSupportCaseStatusOutputBody struct {
	Case *db.SupportCase `json:"case"`
}

type UpdateSupportCaseStatusOutput struct {
	Body UpdateSupportCaseStatusOutputBody `json:"body"`
}

func (h *Handler) UpdateSupportCaseStatus(ctx context.Context, input *UpdateSupportCaseStatusInput) (*UpdateSupportCaseStatusOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	err := h.service.UpdateSupportCaseStatus(ctx, support_cases_service.UpdateSupportCaseStatusInput{
		ID:        input.ID,
		Status:    input.Body.Status,
		UpdatedBy: userID,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update support case status", err)
	}

	supportCase, err := h.service.GetSupportCase(ctx, support_cases_service.GetSupportCaseInput{
		ID: input.ID,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get support case", err)
	}

	return &UpdateSupportCaseStatusOutput{
		Body: UpdateSupportCaseStatusOutputBody{
			Case: supportCase,
		},
	}, nil
}
