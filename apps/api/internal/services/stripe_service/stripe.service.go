package stripe_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/repository/buddy_package_repo"
	"github.com/bored-bots/bored-template/internal/repository/user_repo"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/webhook"
)

type Service struct {
	buddyPackageRepo buddy_package_repo.Repository
	userRepo         user_repo.Repository
	stripeClient     *stripe.Client
	stripeConfig     *config.StripeConfig
}

func New(buddyPackageRepo buddy_package_repo.Repository, userRepo user_repo.Repository, stripeClient *stripe.Client, stripeConfig *config.StripeConfig) *Service {
	return &Service{
		buddyPackageRepo: buddyPackageRepo,
		userRepo:         userRepo,
		stripeConfig:     stripeConfig,
		stripeClient:     stripeClient,
	}
}

func (s *Service) VerifyWebhook(payload []byte, signatureHeader string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signatureHeader, s.stripeConfig.WebhookSecret)
}

func (s *Service) HandleCheckoutSessionCompleted(ctx context.Context, event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("failed to unmarshal checkout session: %w", err)
	}

	if session.Customer == nil {
		return errors.New("session has no customer")
	}

	user, err := s.userRepo.GetUserByStripeCustomerID(ctx, db.GetUserByStripeCustomerIDParams{
		StripeCustomerID: pgtype.Text{
			String: session.Customer.ID,
			Valid:  true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get user by stripe customer id: %w", err)
	}

	retrieve, err := s.stripeClient.V1CheckoutSessions.Retrieve(ctx, session.ID, &stripe.CheckoutSessionRetrieveParams{
		Expand: []*string{stripe.String("line_items")},
	})
	if err != nil {
		return errors.New("failed to retrieve checkout session")
	}

	if retrieve.LineItems == nil || len(retrieve.LineItems.Data) == 0 {
		return errors.New("session has no line items")
	}

	stripePriceID := retrieve.LineItems.Data[0].Price.ID

	packageList, err := s.buddyPackageRepo.GetAllBuddyPackages(ctx)
	if err != nil {
		return fmt.Errorf("failed to get buddy packages: %w", err)
	}

	var purchasedPackage *db.BuddyPackage
	for _, pkg := range packageList {
		if pkg.StripePriceID == stripePriceID {
			purchasedPackage = pkg
			break
		}
	}
	if purchasedPackage == nil {
		return fmt.Errorf("buddy package not found for stripe_price_id: %s", stripePriceID)
	}

	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userRepo := s.userRepo.WithTx(tx)

	err = userRepo.UpdateUserPoints(ctx, db.UpdateUserPointsParams{
		ID:             user.ID,
		LifetimePoints: purchasedPackage.Points,
	})
	if err != nil {
		return fmt.Errorf("failed to update user points: %w", err)
	}

	err = userRepo.InsertBuddyPointsLedger(ctx, db.InsertBuddyPointsLedgerParams{
		UserID:               user.ID,
		TransactionType:      db.TransactionTypePurchase,
		Amount:               purchasedPackage.Points,
		MonthlyPointsBefore:  user.MonthlyPoints,
		LifetimePointsBefore: user.LifetimePoints,
		MonthlyPointsAfter:   user.MonthlyPoints,
		LifetimePointsAfter:  user.LifetimePoints + purchasedPackage.Points,
		Description:          purchasedPackage.Name,
		RelatedEntityID: pgtype.UUID{
			Bytes: purchasedPackage.ID,
			Valid: true,
		},
		RelatedEntityType: db.NullRelatedEntityType{
			RelatedEntityType: db.RelatedEntityTypeBuddyPackage,
			Valid:             true,
		},
	})

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
