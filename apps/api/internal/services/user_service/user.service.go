package user_service

import (
	"context"

	"github.com/bored-bots/bored-template/internal/repository/user_repo"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
)

type Service struct {
	userRepo user_repo.Repository
}

func New(userRepo user_repo.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

type GetBuddyPointsResult struct {
	MonthlyPoints  int32 `json:"monthly_points"`
	LifetimePoints int32 `json:"lifetime_points"`
}

func (s *Service) GetBuddyPoints(ctx context.Context, userID uuid.UUID) (*GetBuddyPointsResult, error) {
	points, err := s.userRepo.GetBuddyPointsByUserID(ctx, db.GetBuddyPointsByUserIDParams{
		ID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &GetBuddyPointsResult{
		MonthlyPoints:  points.MonthlyPoints,
		LifetimePoints: points.LifetimePoints,
	}, nil
}
