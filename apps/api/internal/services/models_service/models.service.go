package models_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidModel        = errors.New("invalid model")
	ErrUnsupportedProvider = errors.New("unsupported model provider")
	ErrModelAlreadyExists  = errors.New("model already exists")
	ErrModelNotFound       = errors.New("model not found")
)

type Repository interface {
	ListEnabledAIModels(ctx context.Context) ([]*db.AiModel, error)
	ListAIModels(ctx context.Context) ([]*db.AiModel, error)
	GetAIModelByProviderModelID(ctx context.Context, providerModelID string) (*db.AiModel, error)
	CreateAIModel(ctx context.Context, params db.CreateAIModelParams) (*db.AiModel, error)
	SetAIModelEnabled(ctx context.Context, params db.SetAIModelEnabledParams) (*db.AiModel, error)
}

type ProviderRegistry interface {
	SupportedProviders() []string
	SupportsProvider(provider string) bool
}

type Service struct {
	repository Repository
	providers  ProviderRegistry
}

func New(repository Repository, providers ProviderRegistry) *Service {
	return &Service{repository: repository, providers: providers}
}

func (s *Service) GetModels(ctx context.Context) ([]*db.AiModel, error) {
	return s.repository.ListAIModels(ctx)
}

func (s *Service) GetEnabledModels(ctx context.Context) ([]*db.AiModel, error) {
	return s.repository.ListEnabledAIModels(ctx)
}

func (s *Service) GetModelByProviderModelID(ctx context.Context, providerModelID string) (*db.AiModel, error) {
	return s.repository.GetAIModelByProviderModelID(ctx, providerModelID)
}

func (s *Service) GetProviders() []string {
	return s.providers.SupportedProviders()
}

type CreateModelInput struct {
	ProviderModelID string
	Name            string
	Provider        string
	SortOrder       int32
}

func (s *Service) CreateModel(ctx context.Context, input CreateModelInput) (*db.AiModel, error) {
	input.ProviderModelID = strings.TrimSpace(input.ProviderModelID)
	input.Name = strings.TrimSpace(input.Name)
	input.Provider = strings.TrimSpace(input.Provider)

	if input.ProviderModelID == "" || input.Name == "" || input.Provider == "" {
		return nil, ErrInvalidModel
	}
	if !s.providers.SupportsProvider(input.Provider) {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, input.Provider)
	}

	model, err := s.repository.CreateAIModel(ctx, db.CreateAIModelParams{
		ProviderModelID: input.ProviderModelID,
		Name:            input.Name,
		Provider:        input.Provider,
		SortOrder:       input.SortOrder,
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, ErrModelAlreadyExists
	}
	return model, err
}

func (s *Service) SetModelEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*db.AiModel, error) {
	model, err := s.repository.SetAIModelEnabled(ctx, db.SetAIModelEnabledParams{
		ID:      id,
		Enabled: enabled,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrModelNotFound
	}
	return model, err
}
