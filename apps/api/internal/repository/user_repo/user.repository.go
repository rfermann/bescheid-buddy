package user_repo

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetBuddyPointsByUserID(ctx context.Context, params db.GetBuddyPointsByUserIDParams) (*db.GetBuddyPointsByUserIDRow, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*db.User, error)
	GetUserByStripeCustomerID(ctx context.Context, params db.GetUserByStripeCustomerIDParams) (*db.User, error)
	UpdateUserPoints(ctx context.Context, params db.UpdateUserPointsParams) error
}

type UserRepository struct {
	Store *store.Store
}

func New(store *store.Store) Repository {
	return &UserRepository{
		Store: store,
	}
}

func (r *UserRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.Store.Pool.Begin(ctx)
}

func (r *UserRepository) GetBuddyPointsByUserID(ctx context.Context, params db.GetBuddyPointsByUserIDParams) (*db.GetBuddyPointsByUserIDRow, error) {
	return r.Store.Queries.GetBuddyPointsByUserID(ctx, params)
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*db.User, error) {
	return r.Store.Queries.GetUserByID(ctx, db.GetUserByIDParams{ID: userID})
}

func (r *UserRepository) GetUserByStripeCustomerID(ctx context.Context, params db.GetUserByStripeCustomerIDParams) (*db.User, error) {
	return r.Store.Queries.GetUserByStripeCustomerID(ctx, db.GetUserByStripeCustomerIDParams{
		StripeCustomerID: params.StripeCustomerID,
	})
}

func (r *UserRepository) UpdateUserPoints(ctx context.Context, params db.UpdateUserPointsParams) error {
	return r.Store.Queries.UpdateUserPoints(ctx, db.UpdateUserPointsParams{
		ID:             params.ID,
		LifetimePoints: params.LifetimePoints,
	})
}
