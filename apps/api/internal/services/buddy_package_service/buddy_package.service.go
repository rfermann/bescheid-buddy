package buddy_package_service

import (
	"context"

	"github.com/bored-bots/bored-template/internal/repository/buddy_package_repo"
	"github.com/bored-bots/bored-template/internal/store/db"
)

type Service struct {
	repository buddy_package_repo.Repository
}

func New(repository buddy_package_repo.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error) {
	return s.repository.GetAllBuddyPackages(ctx)
}
