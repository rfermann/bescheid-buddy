package users_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bored-bots/bored-template/internal/repository/buddy_packages_repository"
	"github.com/bored-bots/bored-template/internal/repository/users_repository"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInsufficientBuddyPoints = errors.New("insufficient BuddyPoints")
	ErrInvalidSpendAmount      = errors.New("spend amount must be positive")
	ErrInvalidRefundAmount     = errors.New("refund amount must be positive")
	ErrBuddyPointSpendNotFound = errors.New("BuddyPoint spend not found")
	ErrBuddyPointSpendMismatch = errors.New("BuddyPoint spend does not match refund")
)

type Service struct {
	userRepo         users_repository.Repository
	buddyPackageRepo buddy_packages_repository.Repository
}

func New(userRepo users_repository.Repository, buddyPackageRepo buddy_packages_repository.Repository) *Service {
	return &Service{
		userRepo:         userRepo,
		buddyPackageRepo: buddyPackageRepo,
	}
}

type GetBuddyPointsResult struct {
	MonthlyPoints  int32 `json:"monthly_points"`
	LifetimePoints int32 `json:"lifetime_points"`
}

type SpendTxInput struct {
	UserID            uuid.UUID
	Amount            int32
	Description       string
	RelatedEntityID   uuid.UUID
	RelatedEntityType db.RelatedEntityType
}

type SpendTxResult struct {
	MonthlyPointsSpent  int32
	LifetimePointsSpent int32
}

type RefundTxInput struct {
	UserID            uuid.UUID
	Amount            int32
	Description       string
	RelatedEntityID   uuid.UUID
	RelatedEntityType db.RelatedEntityType
}

type RefundTxResult struct {
	MonthlyPointsRefunded  int32
	LifetimePointsRefunded int32
}

type LedgerEntry struct {
	ID uuid.UUID `json:"id"`
	// enum values must match db.TransactionType constants in models.go
	TransactionType      db.TransactionType `json:"transaction_type" enum:"purchase,referral,monthly_reset,usage"`
	Amount               int32              `json:"amount"`
	MonthlyPointsBefore  int32              `json:"monthly_points_before"`
	LifetimePointsBefore int32              `json:"lifetime_points_before"`
	MonthlyPointsAfter   int32              `json:"monthly_points_after"`
	LifetimePointsAfter  int32              `json:"lifetime_points_after"`
	Description          string             `json:"description"`
	RelatedEntityID      *string            `json:"related_entity_id,omitempty"`
	RelatedEntityType    *string            `json:"related_entity_type,omitempty"`
	RelatedEntityName    *string            `json:"related_entity_name,omitempty"`
	CreatedAt            string             `json:"created_at"`
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

func (s *Service) SpendTx(ctx context.Context, tx pgx.Tx, input SpendTxInput) (*SpendTxResult, error) {
	if input.Amount <= 0 {
		return nil, ErrInvalidSpendAmount
	}

	userRepo := s.userRepo.WithTx(tx)
	points, err := userRepo.GetBuddyPointsForUpdate(ctx, db.GetBuddyPointsForUpdateParams{
		UserID: input.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("lock BuddyPoint balance: %w", err)
	}
	if points.MonthlyPoints < 0 || points.LifetimePoints < 0 {
		return nil, fmt.Errorf("invalid negative BuddyPoint balance for user %s", input.UserID)
	}

	if int64(points.MonthlyPoints)+int64(points.LifetimePoints) < int64(input.Amount) {
		return nil, ErrInsufficientBuddyPoints
	}

	monthlySpent := min(points.MonthlyPoints, input.Amount)
	lifetimeSpent := input.Amount - monthlySpent
	monthlyAfter := points.MonthlyPoints - monthlySpent
	lifetimeAfter := points.LifetimePoints - lifetimeSpent

	rowsAffected, err := userRepo.DebitBuddyPoints(ctx, db.DebitBuddyPointsParams{
		UserID:         input.UserID,
		MonthlyPoints:  monthlySpent,
		LifetimePoints: lifetimeSpent,
	})
	if err != nil {
		return nil, fmt.Errorf("debit BuddyPoints: %w", err)
	}
	if rowsAffected != 1 {
		return nil, ErrInsufficientBuddyPoints
	}

	if err := userRepo.InsertBuddyPointsLedger(ctx, db.InsertBuddyPointsLedgerParams{
		UserID:               input.UserID,
		TransactionType:      db.TransactionTypeUsage,
		Amount:               -input.Amount,
		MonthlyPointsBefore:  points.MonthlyPoints,
		LifetimePointsBefore: points.LifetimePoints,
		MonthlyPointsAfter:   monthlyAfter,
		LifetimePointsAfter:  lifetimeAfter,
		Description:          input.Description,
		RelatedEntityID: pgtype.UUID{
			Bytes: input.RelatedEntityID,
			Valid: true,
		},
		RelatedEntityType: db.NullRelatedEntityType{
			RelatedEntityType: input.RelatedEntityType,
			Valid:             true,
		},
	}); err != nil {
		return nil, fmt.Errorf("insert BuddyPoint spend ledger entry: %w", err)
	}

	return &SpendTxResult{
		MonthlyPointsSpent:  monthlySpent,
		LifetimePointsSpent: lifetimeSpent,
	}, nil
}

func (s *Service) RefundTx(ctx context.Context, tx pgx.Tx, input RefundTxInput) (*RefundTxResult, error) {
	if input.Amount <= 0 {
		return nil, ErrInvalidRefundAmount
	}

	userRepo := s.userRepo.WithTx(tx)
	spend, err := userRepo.GetMatchingBuddyPointSpend(ctx, db.GetMatchingBuddyPointSpendParams{
		UserID: input.UserID,
		RelatedEntityID: pgtype.UUID{
			Bytes: input.RelatedEntityID,
			Valid: true,
		},
		RelatedEntityType: db.NullRelatedEntityType{
			RelatedEntityType: input.RelatedEntityType,
			Valid:             true,
		},
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrBuddyPointSpendNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load BuddyPoint spend: %w", err)
	}

	monthlyRefund := int64(spend.MonthlyPointsBefore) - int64(spend.MonthlyPointsAfter)
	lifetimeRefund := int64(spend.LifetimePointsBefore) - int64(spend.LifetimePointsAfter)
	totalRefund := monthlyRefund + lifetimeRefund
	if monthlyRefund < 0 || lifetimeRefund < 0 ||
		totalRefund != -int64(spend.Amount) || totalRefund != int64(input.Amount) {
		return nil, ErrBuddyPointSpendMismatch
	}

	points, err := userRepo.GetBuddyPointsForUpdate(ctx, db.GetBuddyPointsForUpdateParams{
		UserID: input.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("lock BuddyPoint balance: %w", err)
	}
	if points.MonthlyPoints < 0 || points.LifetimePoints < 0 {
		return nil, fmt.Errorf("invalid negative BuddyPoint balance for user %s", input.UserID)
	}

	monthlyAfter := int64(points.MonthlyPoints) + monthlyRefund
	lifetimeAfter := int64(points.LifetimePoints) + lifetimeRefund
	if monthlyAfter > int64(^uint32(0)>>1) || lifetimeAfter > int64(^uint32(0)>>1) {
		return nil, fmt.Errorf("BuddyPoint balance overflow for user %s", input.UserID)
	}

	monthlyRefunded := int32(monthlyRefund)
	lifetimeRefunded := int32(lifetimeRefund)
	rowsAffected, err := userRepo.CreditBuddyPoints(ctx, db.CreditBuddyPointsParams{
		UserID:         input.UserID,
		MonthlyPoints:  monthlyRefunded,
		LifetimePoints: lifetimeRefunded,
	})
	if err != nil {
		return nil, fmt.Errorf("credit BuddyPoints: %w", err)
	}
	if rowsAffected != 1 {
		return nil, fmt.Errorf("credit BuddyPoints: user %s not found", input.UserID)
	}

	if err := userRepo.InsertBuddyPointsLedger(ctx, db.InsertBuddyPointsLedgerParams{
		UserID:               input.UserID,
		TransactionType:      db.TransactionTypeUsage,
		Amount:               input.Amount,
		MonthlyPointsBefore:  points.MonthlyPoints,
		LifetimePointsBefore: points.LifetimePoints,
		MonthlyPointsAfter:   int32(monthlyAfter),
		LifetimePointsAfter:  int32(lifetimeAfter),
		Description:          input.Description,
		RelatedEntityID: pgtype.UUID{
			Bytes: input.RelatedEntityID,
			Valid: true,
		},
		RelatedEntityType: db.NullRelatedEntityType{
			RelatedEntityType: input.RelatedEntityType,
			Valid:             true,
		},
	}); err != nil {
		return nil, fmt.Errorf("insert BuddyPoint refund ledger entry: %w", err)
	}

	return &RefundTxResult{
		MonthlyPointsRefunded:  monthlyRefunded,
		LifetimePointsRefunded: lifetimeRefunded,
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
			TransactionType:      entry.TransactionType,
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

			if entry.RelatedEntityType.RelatedEntityType == db.RelatedEntityTypeBuddyPackages && entry.RelatedEntityID.Valid {
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
