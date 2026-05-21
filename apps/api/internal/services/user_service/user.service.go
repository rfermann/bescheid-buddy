package user_service

import (
	"context"
	"time"

	"github.com/bored-bots/bored-template/internal/repository/buddy_package_repo"
	"github.com/bored-bots/bored-template/internal/repository/user_repo"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
)

type Service struct {
	userRepo         user_repo.Repository
	buddyPackageRepo buddy_package_repo.Repository
}

func New(userRepo user_repo.Repository, buddyPackageRepo buddy_package_repo.Repository) *Service {
	return &Service{
		userRepo:         userRepo,
		buddyPackageRepo: buddyPackageRepo,
	}
}

type GetBuddyPointsResult struct {
	MonthlyPoints  int32 `json:"monthly_points"`
	LifetimePoints int32 `json:"lifetime_points"`
}

type LedgerEntry struct {
	ID                   uuid.UUID `json:"id"`
	TransactionType      string    `json:"transaction_type"`
	Amount               int32     `json:"amount"`
	MonthlyPointsBefore  int32     `json:"monthly_points_before"`
	LifetimePointsBefore int32     `json:"lifetime_points_before"`
	MonthlyPointsAfter   int32     `json:"monthly_points_after"`
	LifetimePointsAfter  int32     `json:"lifetime_points_after"`
	Description          string    `json:"description"`
	RelatedEntityID      *string   `json:"related_entity_id,omitempty"`
	RelatedEntityType    *string   `json:"related_entity_type,omitempty"`
	RelatedEntityName    *string   `json:"related_entity_name,omitempty"`
	CreatedAt            string    `json:"created_at"`
}

type GetBuddyPointsLedgerResult struct {
	Entries []LedgerEntry `json:"entries"`
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

func (s *Service) GetBuddyPointsLedger(ctx context.Context, userID uuid.UUID) (*GetBuddyPointsLedgerResult, error) {
	ledger, err := s.userRepo.GetBuddyPointsLedgerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	entries := make([]LedgerEntry, 0, len(ledger))
	for _, entry := range ledger {
		ledgerEntry := LedgerEntry{
			ID:                   entry.ID,
			TransactionType:      string(entry.TransactionType),
			Amount:               entry.Amount,
			MonthlyPointsBefore:  entry.MonthlyPointsBefore,
			LifetimePointsBefore: entry.LifetimePointsBefore,
			MonthlyPointsAfter:   entry.MonthlyPointsAfter,
			LifetimePointsAfter:  entry.LifetimePointsAfter,
			Description:          entry.Description,
			CreatedAt:            entry.CreatedAt.Format(time.RFC3339),
		}

		if entry.RelatedEntityID.Valid {
			relatedEntityUUID, _ := uuid.FromBytes(entry.RelatedEntityID.Bytes[:])
			relatedEntityID := relatedEntityUUID.String()
			ledgerEntry.RelatedEntityID = &relatedEntityID
		}

		if entry.RelatedEntityType.Valid {
			relatedEntityType := string(entry.RelatedEntityType.RelatedEntityType)
			ledgerEntry.RelatedEntityType = &relatedEntityType

			if entry.RelatedEntityType.RelatedEntityType == db.RelatedEntityTypeBuddyPackage && entry.RelatedEntityID.Valid {
				relatedEntityUUID, _ := uuid.FromBytes(entry.RelatedEntityID.Bytes[:])
				buddyPackage, err := s.buddyPackageRepo.GetBuddyPackageByID(ctx, db.GetBuddyPackageByIDParams{
					ID: relatedEntityUUID,
				})
				if err == nil && buddyPackage != nil {
					ledgerEntry.RelatedEntityName = &buddyPackage.Name
				}
			}
		}

		entries = append(entries, ledgerEntry)
	}

	return &GetBuddyPointsLedgerResult{
		Entries: entries,
	}, nil
}
