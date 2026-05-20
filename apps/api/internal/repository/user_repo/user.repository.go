package user_repo

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
)

type Repository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*db.User, error)
	GetUserByStripeCustomerID(ctx context.Context, params db.GetUserByStripeCustomerIDParams) (*db.User, error)
	UpdateUserPoints(ctx context.Context, params db.UpdateUserPointsParams) error
}

type userRepository struct {
	store *store.Store
}

func New(store *store.Store) Repository {
	return &userRepository{
		store: store,
	}
}

func (r *userRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*db.User, error) {
	return r.store.Queries.GetUserByID(ctx, db.GetUserByIDParams{ID: userID})
}

func (r *userRepository) GetUserByStripeCustomerID(ctx context.Context, params db.GetUserByStripeCustomerIDParams) (*db.User, error) {
	return r.store.Queries.GetUserByStripeCustomerID(ctx, db.GetUserByStripeCustomerIDParams{
		StripeCustomerID: params.StripeCustomerID,
	})
}

func (r *userRepository) UpdateUserPoints(ctx context.Context, params db.UpdateUserPointsParams) error {
	return r.store.Queries.UpdateUserPoints(ctx, db.UpdateUserPointsParams{
		ID:             params.ID,
		LifetimePoints: params.LifetimePoints,
	})
}
