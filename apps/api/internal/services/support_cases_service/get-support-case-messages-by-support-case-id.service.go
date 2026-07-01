package support_cases_service

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
)

type GetSupportCaseMessagesBySupportCaseIDInput struct {
	CaseID uuid.UUID
	UserID uuid.UUID
}

func (s *Service) GetSupportCaseMessagesBySupportCaseID(ctx context.Context, input *GetSupportCaseMessagesBySupportCaseIDInput) ([]*db.SupportCaseMessage, error) {
	return s.repository.GetSupportCaseMessagesBySupportCaseID(ctx, db.GetSupportCaseMessagesBySupportCaseIDParams{
		CaseID: input.CaseID,
		UserID: input.UserID,
	})
}
