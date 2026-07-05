package support_cases_handler

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
)

type GetAllSupportCasesOutputBody struct {
	Cases []*db.SupportCase `json:"cases"`
}

type GetAllSupportCasesOutput struct {
	Body GetAllSupportCasesOutputBody `json:"body"`
}

func (h *Handler) GetAllSupportCases(ctx context.Context, _ *struct{}) (*GetAllSupportCasesOutput, error) {
	supportCases, err := h.service.GetAllSupportCases(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get support cases", err)
	}

	return &GetAllSupportCasesOutput{
		Body: GetAllSupportCasesOutputBody{
			Cases: supportCases,
		},
	}, nil
}
