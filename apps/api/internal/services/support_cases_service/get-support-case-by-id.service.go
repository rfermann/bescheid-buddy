package support_cases_service

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
)

type GetSupportCaseByIDInput struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (s *Service) GetSupportCaseByID(ctx context.Context, params GetSupportCaseByIDInput) (*db.SupportCase, error) {
	return s.repository.GetSupportCaseByID(ctx, db.GetSupportCaseByIDParams{
		ID:     params.ID,
		UserID: params.UserID,
	})
}
