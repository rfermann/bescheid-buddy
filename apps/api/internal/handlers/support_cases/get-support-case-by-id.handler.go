package support_cases_handler

import (
	"context"
	"fmt"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetSupportCaseByIDInput struct {
	ID uuid.UUID `path:"id"`
}

type GetSupportCaseByIDOutputBody struct {
	SupportCase *db.SupportCase `json:"supportCase"`
}

type GetSupportCaseByIDOutput struct {
	Body GetSupportCaseByIDOutputBody `json:"body"`
}

func (h *Handler) GetSupportCaseByID(ctx context.Context, input *GetSupportCaseByIDInput) (*GetSupportCaseByIDOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	supportCase, err := h.service.GetSupportCaseByID(ctx, support_cases_service.GetSupportCaseByIDInput{
		ID:     input.ID,
		UserID: userID,
	})
	if err != nil {
		fmt.Println(err)
		return nil, huma.Error404NotFound("support case not found")
	}

	return &GetSupportCaseByIDOutput{
		Body: GetSupportCaseByIDOutputBody{
			SupportCase: supportCase,
		},
	}, nil
}
