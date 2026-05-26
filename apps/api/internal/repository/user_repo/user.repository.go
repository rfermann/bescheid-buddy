package user_repo

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	GetBuddyPointsByUserID(ctx context.Context, params db.GetBuddyPointsByUserIDParams) (*db.GetBuddyPointsByUserIDRow, error)
	GetBuddyPointsLedgerByUserID(ctx context.Context, userID uuid.UUID) ([]*db.BuddyPointsLedger, error)
	GetUserByID(ctx context.Context, userID db.GetUserByIDParams) (*db.User, error)
	GetUserByStripeCustomerID(ctx context.Context, params db.GetUserByStripeCustomerIDParams) (*db.User, error)
	InsertBuddyPointsLedger(ctx context.Context, params db.InsertBuddyPointsLedgerParams) error
	UpdateUserPoints(ctx context.Context, params db.UpdateUserPointsParams) error
	UpdateUserStripeCustomerID(ctx context.Context, params db.UpdateUserStripeCustomerIDParams) error
	WithTx(tx pgx.Tx) Repository
}

type UserRepository struct {
	Store *store.Store
}

func New(store *store.Store) Repository {
	return &UserRepository{
		Store: store,
	}
}

func (r *UserRepository) WithTx(tx pgx.Tx) Repository {
	return &UserRepository{
		Store: &store.Store{
			Pool:    r.Store.Pool,
			Queries: db.New(tx),
		},
	}
}

func (r *UserRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.Store.Pool.Begin(ctx)
}

func (r *UserRepository) GetBuddyPointsByUserID(ctx context.Context, params db.GetBuddyPointsByUserIDParams) (*db.GetBuddyPointsByUserIDRow, error) {
	return r.Store.Queries.GetBuddyPointsByUserID(ctx, params)
}

func (r *UserRepository) GetBuddyPointsLedgerByUserID(ctx context.Context, userID uuid.UUID) ([]*db.BuddyPointsLedger, error) {
	return r.Store.Queries.GetBuddyPointsLedgerByUserID(ctx, db.GetBuddyPointsLedgerByUserIDParams{
		UserID: userID,
	})
}

func (r *UserRepository) GetUserByID(ctx context.Context, params db.GetUserByIDParams) (*db.User, error) {
	return r.Store.Queries.GetUserByID(ctx, params)
}

func (r *UserRepository) GetUserByStripeCustomerID(ctx context.Context, params db.GetUserByStripeCustomerIDParams) (*db.User, error) {
	return r.Store.Queries.GetUserByStripeCustomerID(ctx, db.GetUserByStripeCustomerIDParams{
		StripeCustomerID: params.StripeCustomerID,
	})
}

func (r *UserRepository) InsertBuddyPointsLedger(ctx context.Context, params db.InsertBuddyPointsLedgerParams) error {
	return r.Store.Queries.InsertBuddyPointsLedger(ctx, params)
}

func (r *UserRepository) UpdateUserPoints(ctx context.Context, params db.UpdateUserPointsParams) error {
	return r.Store.Queries.UpdateUserPoints(ctx, db.UpdateUserPointsParams{
		ID:             params.ID,
		LifetimePoints: params.LifetimePoints,
	})
}

func (r *UserRepository) UpdateUserStripeCustomerID(ctx context.Context, params db.UpdateUserStripeCustomerIDParams) error {
	return r.Store.Queries.UpdateUserStripeCustomerID(ctx, db.UpdateUserStripeCustomerIDParams{
		ID:               params.ID,
		StripeCustomerID: params.StripeCustomerID,
	})
}
