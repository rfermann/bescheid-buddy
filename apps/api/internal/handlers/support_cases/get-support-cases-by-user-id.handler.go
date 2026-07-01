package support_cases_handler

import (
	"context"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
)

type GetSupportCasesByUserIdOutputBody struct {
	Cases []*db.SupportCase `json:"cases"`
}

type GetSupportCasesByUserIdOutput struct {
	Body GetSupportCasesByUserIdOutputBody `json:"body"`
}

func (h *Handler) GetSupportCasesByUserId(ctx context.Context, _ *struct{}) (*GetSupportCasesByUserIdOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	supportCases, err := h.service.GetSupportCasesByUserId(ctx, db.GetSupportCasesByUserIDParams{
		UserID: userID,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get support cases", err)
	}

	return &GetSupportCasesByUserIdOutput{
		Body: GetSupportCasesByUserIdOutputBody{
			Cases: supportCases,
		},
	}, nil
}
