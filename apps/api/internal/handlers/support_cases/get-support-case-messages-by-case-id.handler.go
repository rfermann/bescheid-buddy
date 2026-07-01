package support_cases_handler

import (
	"context"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetSupportCaseMessagesBySupportCaseIDInput struct {
	CaseID uuid.UUID `path:"caseId"`
}

type GetSupportCaseMessagesBySupportCaseIDOutputBody struct {
	Messages []*db.SupportCaseMessage `json:"messages"`
}

type GetSupportCaseMessagesBySupportCaseIDOutput struct {
	Body GetSupportCaseMessagesBySupportCaseIDOutputBody `json:"body"`
}

func (h *Handler) GetSupportCaseMessagesBySupportCaseID(ctx context.Context, input *GetSupportCaseMessagesBySupportCaseIDInput) (*GetSupportCaseMessagesBySupportCaseIDOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	messages, err := h.service.GetSupportCaseMessagesBySupportCaseID(ctx, &support_cases_service.GetSupportCaseMessagesBySupportCaseIDInput{
		CaseID: input.CaseID,
		UserID: userID,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get support case messages", err)
	}

	return &GetSupportCaseMessagesBySupportCaseIDOutput{
		Body: GetSupportCaseMessagesBySupportCaseIDOutputBody{
			Messages: messages,
		},
	}, nil
}
