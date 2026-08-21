package models_repository

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
)

type Repository interface {
	ListEnabledAIModels(ctx context.Context) ([]*db.AiModel, error)
	ListAIModels(ctx context.Context) ([]*db.AiModel, error)
	GetAIModelByProviderModelID(ctx context.Context, providerModelID string) (*db.AiModel, error)
	CreateAIModel(ctx context.Context, params db.CreateAIModelParams) (*db.AiModel, error)
	SetAIModelEnabled(ctx context.Context, params db.SetAIModelEnabledParams) (*db.AiModel, error)
}

type modelsRepository struct {
	store *store.Store
}

func New(store *store.Store) Repository {
	return &modelsRepository{store: store}
}

func (r *modelsRepository) ListEnabledAIModels(ctx context.Context) ([]*db.AiModel, error) {
	return r.store.Queries.ListEnabledAIModels(ctx)
}

func (r *modelsRepository) ListAIModels(ctx context.Context) ([]*db.AiModel, error) {
	return r.store.Queries.ListAIModels(ctx)
}

func (r *modelsRepository) GetAIModelByProviderModelID(ctx context.Context, providerModelID string) (*db.AiModel, error) {
	return r.store.Queries.GetAIModelByProviderModelID(ctx, db.GetAIModelByProviderModelIDParams{
		ProviderModelID: providerModelID,
	})
}

func (r *modelsRepository) CreateAIModel(ctx context.Context, params db.CreateAIModelParams) (*db.AiModel, error) {
	return r.store.Queries.CreateAIModel(ctx, params)
}

func (r *modelsRepository) SetAIModelEnabled(ctx context.Context, params db.SetAIModelEnabledParams) (*db.AiModel, error) {
	return r.store.Queries.SetAIModelEnabled(ctx, params)
}
