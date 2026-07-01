package support_cases_service

import (
	"context"
	"fmt"

	"github.com/bored-bots/bored-template/internal/repository/support_cases_repository"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type Service struct {
	repository support_cases_repository.Repository
}

func New(repository support_cases_repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

type CreateSupportCaseInput struct {
	Category db.SupportCaseCategory
	Message  string
	Subject  string
	UserID   uuid.UUID
}

func (s *Service) CreateSupportCase(ctx context.Context, params CreateSupportCaseInput) (*db.SupportCase, error) {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	supportCaseRepository := s.repository.WithTx(tx)
	supportCase, err := supportCaseRepository.CreateSupportCase(ctx, db.CreateSupportCaseParams{
		UserID:    params.UserID,
		Subject:   params.Subject,
		Status:    db.SupportCaseStatusOpen,
		Category:  params.Category,
		UpdatedBy: params.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create support case: %w", err)
	}

	_, err = supportCaseRepository.CreateSupportCaseMessage(ctx, db.CreateSupportCaseMessageParams{
		CaseID:        supportCase.ID,
		UserID:        params.UserID,
		Message:       params.Message,
		IsUserMessage: true,
		IsInternal:    false,
		UpdatedBy:     params.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create support case message: %w", err)
	}
	tx.Commit(ctx)

	return supportCase, nil
}

func (s *Service) GetSupportCasesByUserId(ctx context.Context, input db.GetSupportCasesByUserIDParams) ([]*db.SupportCase, error) {
	return s.repository.GetSupportCasesByUserID(ctx, input)
}
