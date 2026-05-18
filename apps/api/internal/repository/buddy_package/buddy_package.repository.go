package buddy_package

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
)

type BuddyPackageRepository interface {
	GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error)
}

type buddyPackageRepository struct {
	store *store.Store
}

func NewBuddyPackageRepository(store *store.Store) BuddyPackageRepository {
	return &buddyPackageRepository{
		store: store,
	}
}

func (r *buddyPackageRepository) GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error) {
	return r.store.Queries.GetAllBuddyPackages(ctx)
}
