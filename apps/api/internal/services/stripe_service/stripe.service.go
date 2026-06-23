package stripe_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/repository/buddy_packages_repository"
	"github.com/bored-bots/bored-template/internal/repository/users_repository"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/webhook"
)

type Service struct {
	buddyPackageRepo buddy_packages_repository.Repository
	userRepo         users_repository.Repository
	stripeClient     *stripe.Client
	stripeConfig     *config.StripeConfig
}

func New(buddyPackageRepo buddy_packages_repository.Repository, userRepo users_repository.Repository, stripeClient *stripe.Client, stripeConfig *config.StripeConfig) *Service {
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
			RelatedEntityType: db.RelatedEntityTypeBuddyPackages,
			Valid:             true,
		},
	})

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type CreateCheckoutSessionInput struct {
	UserID        uuid.UUID
	StripePriceID string
	PackageName   string
}

type CreateCheckoutSessionResult struct {
	URL string
}

func (s *Service) CreateCheckoutSession(ctx context.Context, input CreateCheckoutSessionInput) (*CreateCheckoutSessionResult, error) {
	user, err := s.userRepo.GetUserByID(ctx, db.GetUserByIDParams{
		ID: input.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	stripeCustomerID := user.StripeCustomerID.String
	if !user.StripeCustomerID.Valid || stripeCustomerID == "" {
		customer, err := s.stripeClient.V1Customers.Create(ctx, &stripe.CustomerCreateParams{
			Name:  stripe.String(user.Name),
			Email: stripe.String(user.Email),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create stripe customer: %w", err)
		}

		stripeCustomerID = customer.ID

		err = s.userRepo.UpdateUserStripeCustomerID(ctx, db.UpdateUserStripeCustomerIDParams{
			ID:               user.ID,
			StripeCustomerID: pgtype.Text{String: stripeCustomerID, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update user stripe customer id: %w", err)
		}
	}

	params := &stripe.CheckoutSessionCreateParams{
		Customer: stripe.String(stripeCustomerID),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				Price:    stripe.String(input.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s/buddypoints-verwalten?package_name=%s", s.stripeConfig.Origin, input.PackageName)),
		CancelURL:  stripe.String(fmt.Sprintf("%s/buddypoints-verwalten?canceled=true", s.stripeConfig.Origin)),
	}

	session, err := s.stripeClient.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return &CreateCheckoutSessionResult{
		URL: session.URL,
	}, nil
}
