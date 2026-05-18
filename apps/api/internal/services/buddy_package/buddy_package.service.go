package buddy_package

import (
	"context"

	"github.com/bored-bots/bored-template/internal/repository"
	"github.com/bored-bots/bored-template/internal/store/db"
)

type Service struct {
	repository repository.BuddyPackageRepository
}

func NewService(repository repository.BuddyPackageRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error) {
	return s.repository.GetAllBuddyPackages(ctx)
}
