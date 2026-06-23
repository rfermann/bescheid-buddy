package buddy_packages_service

import (
	"context"

	"github.com/bored-bots/bored-template/internal/repository/buddy_packages_repository"
	"github.com/bored-bots/bored-template/internal/store/db"
)

type Service struct {
	repository buddy_packages_repository.Repository
}

func New(repository buddy_packages_repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllBuddyPackages(ctx context.Context) ([]*db.BuddyPackage, error) {
	return s.repository.GetAllBuddyPackages(ctx)
}
