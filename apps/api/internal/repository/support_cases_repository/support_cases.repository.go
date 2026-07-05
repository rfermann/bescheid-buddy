package support_cases_repository

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreateSupportCase(ctx context.Context, params db.CreateSupportCaseParams) (*db.SupportCase, error)
	CreateSupportCaseMessage(ctx context.Context, params db.CreateSupportCaseMessageParams) (*db.SupportCaseMessage, error)
	GetSupportCaseByID(ctx context.Context, params db.GetSupportCaseByIDParams) (*db.SupportCase, error)
	GetSupportCase(ctx context.Context, params db.GetSupportCaseParams) (*db.SupportCase, error)
	GetSupportCasesByUserID(ctx context.Context, params db.GetSupportCasesByUserIDParams) ([]*db.SupportCase, error)
	GetSupportCaseMessagesBySupportCaseID(ctx context.Context, params db.GetSupportCaseMessagesBySupportCaseIDParams) ([]*db.SupportCaseMessage, error)
	GetAllSupportCases(ctx context.Context) ([]*db.SupportCase, error)
	UpdateSupportCaseStatus(ctx context.Context, params db.UpdateSupportCaseStatusParams) error
	WithTx(tx pgx.Tx) Repository
}

type SupportCasesRepository struct {
	Store *store.Store
}

func New(store *store.Store) Repository {
	return &SupportCasesRepository{
		Store: store,
	}
}

func (r *SupportCasesRepository) WithTx(tx pgx.Tx) Repository {
	return &SupportCasesRepository{
		Store: &store.Store{
			Pool:    r.Store.Pool,
			Queries: db.New(tx),
		},
	}
}

func (r *SupportCasesRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.Store.Pool.Begin(ctx)
}

func (r *SupportCasesRepository) CreateSupportCase(ctx context.Context, params db.CreateSupportCaseParams) (*db.SupportCase, error) {
	return r.Store.Queries.CreateSupportCase(ctx, params)
}

func (r *SupportCasesRepository) CreateSupportCaseMessage(ctx context.Context, params db.CreateSupportCaseMessageParams) (*db.SupportCaseMessage, error) {
	return r.Store.Queries.CreateSupportCaseMessage(ctx, params)
}

func (r *SupportCasesRepository) GetSupportCaseByID(ctx context.Context, params db.GetSupportCaseByIDParams) (*db.SupportCase, error) {
	return r.Store.Queries.GetSupportCaseByID(ctx, params)
}

func (r *SupportCasesRepository) GetSupportCase(ctx context.Context, params db.GetSupportCaseParams) (*db.SupportCase, error) {
	return r.Store.Queries.GetSupportCase(ctx, params)
}

func (r *SupportCasesRepository) GetSupportCasesByUserID(ctx context.Context, params db.GetSupportCasesByUserIDParams) ([]*db.SupportCase, error) {
	return r.Store.Queries.GetSupportCasesByUserID(ctx, params)
}

func (r *SupportCasesRepository) GetSupportCaseMessagesBySupportCaseID(ctx context.Context, params db.GetSupportCaseMessagesBySupportCaseIDParams) ([]*db.SupportCaseMessage, error) {
	return r.Store.Queries.GetSupportCaseMessagesBySupportCaseID(ctx, params)
}

func (r *SupportCasesRepository) GetAllSupportCases(ctx context.Context) ([]*db.SupportCase, error) {
	return r.Store.Queries.GetAllSupportCases(ctx)
}

func (r *SupportCasesRepository) UpdateSupportCaseStatus(ctx context.Context, params db.UpdateSupportCaseStatusParams) error {
	return r.Store.Queries.UpdateSupportCaseStatus(ctx, params)
}
