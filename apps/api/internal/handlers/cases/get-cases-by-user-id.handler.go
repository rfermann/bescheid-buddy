package cases_handler

import (
	"context"
	"time"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CaseListItem struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetCasesByUserIDOutput struct {
	Body struct {
		Cases []CaseListItem `json:"cases"`
	}
}

func (h *Handler) getCasesByUserID(ctx context.Context, _ *struct{}) (*GetCasesByUserIDOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	cases, err := h.service.GetCasesByUserID(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get cases", err)
	}

	output := &GetCasesByUserIDOutput{}
	output.Body.Cases = make([]CaseListItem, len(cases))
	for i, c := range cases {
		output.Body.Cases[i] = CaseListItem{
			ID:        c.ID,
			Title:     c.Title.String,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	return output, nil
}
