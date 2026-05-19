package buddy_package

import (
	"context"

	"github.com/bored-bots/bored-template/internal/repository/buddy_package"
	"github.com/bored-bots/bored-template/internal/store/db"
)

type Service struct {
	repository buddy_package.BuddyPackageRepository
}

func NewService(repository buddy_package.BuddyPackageRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error) {
	return s.repository.GetAllBuddyPackages(ctx)
}
