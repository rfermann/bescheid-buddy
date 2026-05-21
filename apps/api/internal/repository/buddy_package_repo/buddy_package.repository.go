package buddy_package_repo

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
)

type Repository interface {
	GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error)
	GetBuddyPackageByID(ctx context.Context, params db.GetBuddyPackageByIDParams) (*db.BuddyPackage, error)
}

type buddyPackageRepository struct {
	store *store.Store
}

func New(store *store.Store) Repository {
	return &buddyPackageRepository{
		store: store,
	}
}

func (r *buddyPackageRepository) GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error) {
	return r.store.Queries.GetAllBuddyPackages(ctx)
}

func (r *buddyPackageRepository) GetBuddyPackageByID(ctx context.Context, params db.GetBuddyPackageByIDParams) (*db.BuddyPackage, error) {
	return r.store.Queries.GetBuddyPackageByID(ctx, params)
}
