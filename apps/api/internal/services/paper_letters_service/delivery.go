package paper_letters_service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/services/users_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrDeliveryNotFound   = errors.New("paper letter delivery was not found")
	ErrUploadNotConfirmed = errors.New("paper letter upload is not confirmed")
)

const deliveryRequestConstraint = "paper_letter_deliveries_user_client_request_unique"

type DeliveryProvider interface {
	Submit(ctx context.Context, pdf []byte, filename, idempotencyKey, product, addressPosition, printMode, printSpectrum string) (string, error)
}

type PointRefunder interface {
	RefundTx(context.Context, pgx.Tx, users_service.RefundTxInput) (*users_service.RefundTxResult, error)
}

type DeliveryProcessor struct {
	repository  cases_repository.Repository
	storage     PaperLetterStorage
	provider    DeliveryProvider
	refunder    PointRefunder
	maxPDFBytes int64
	now         func() time.Time
}

func NewDeliveryProcessor(repository cases_repository.Repository, storage PaperLetterStorage, provider DeliveryProvider, refunder PointRefunder, maxPDFBytes int64) *DeliveryProcessor {
	return &DeliveryProcessor{repository: repository, storage: storage, provider: provider, refunder: refunder, maxPDFBytes: maxPDFBytes, now: time.Now}
}

type CreateDeliveryInput struct {
	UserID          uuid.UUID
	UploadID        uuid.UUID
	ClientRequestID uuid.UUID
}

type CreateDeliveryResult struct {
	Delivery *db.PaperLetterDelivery
	Existing bool
}

func (s *Service) CreateDelivery(ctx context.Context, input CreateDeliveryInput) (*CreateDeliveryResult, error) {
	if input.ClientRequestID == uuid.Nil {
		return nil, ErrInvalidInput
	}
	if existing, err := s.deliveryByRequest(ctx, input.UserID, input.ClientRequestID); err == nil {
		return &CreateDeliveryResult{Delivery: existing, Existing: true}, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("load delivery request: %w", err)
	}

	deliveryID := uuid.New()
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin delivery transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	txRepo := s.repository.WithTx(tx)
	upload, err := txRepo.GetPaperLetterUploadByIDAndUserID(ctx, db.GetPaperLetterUploadByIDAndUserIDParams{ID: input.UploadID, UserID: input.UserID})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPaperLetterUploadNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load delivery upload: %w", err)
	}
	if upload.Status != db.PaperLetterUploadStatusConfirmed {
		return nil, ErrUploadNotConfirmed
	}
	cost := int32(deliveryPrice(int(upload.PageCount)))
	if _, err := s.points.SpendTx(ctx, tx, users_service.SpendTxInput{
		UserID: input.UserID, Amount: cost, Description: "Papierbrief delivery",
		RelatedEntityID: deliveryID, RelatedEntityType: db.RelatedEntityTypePaperLetterDelivery,
	}); err != nil {
		return nil, fmt.Errorf("spend delivery BuddyPoints: %w", err)
	}
	delivery, err := txRepo.CreatePaperLetterDelivery(ctx, db.CreatePaperLetterDeliveryParams{
		ID: deliveryID, UserID: input.UserID, UploadID: input.UploadID, ClientRequestID: input.ClientRequestID,
		PageCount: upload.PageCount, BpCost: cost, DeliveryProduct: s.pingen.Product,
		AddressPosition: s.pingen.AddressPosition, PrintMode: s.pingen.PrintMode, PrintSpectrum: s.pingen.PrintSpectrum,
	})
	if err != nil {
		if isDeliveryRequestConflict(err) {
			_ = tx.Rollback(ctx)
			existing, loadErr := s.deliveryByRequest(ctx, input.UserID, input.ClientRequestID)
			if loadErr != nil {
				return nil, fmt.Errorf("load concurrent delivery request: %w", loadErr)
			}
			return &CreateDeliveryResult{Delivery: existing, Existing: true}, nil
		}
		return nil, fmt.Errorf("create delivery: %w", err)
	}
	jobID, err := s.queue.QueuePaperLetterDeliveryJobTx(ctx, tx, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("enqueue delivery: %w", err)
	}
	if err := txRepo.SetPaperLetterDeliveryJobID(ctx, db.SetPaperLetterDeliveryJobIDParams{ID: deliveryID, JobID: pgtype.Int8{Int64: jobID, Valid: true}}); err != nil {
		return nil, fmt.Errorf("store delivery job ID: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit delivery: %w", err)
	}
	delivery.JobID = pgtype.Int8{Int64: jobID, Valid: true}
	return &CreateDeliveryResult{Delivery: delivery}, nil
}

func (s *Service) GetDelivery(ctx context.Context, userID, deliveryID uuid.UUID) (*db.PaperLetterDelivery, error) {
	delivery, err := s.repository.GetPaperLetterDeliveryByID(ctx, db.GetPaperLetterDeliveryByIDParams{ID: deliveryID})
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && delivery.UserID != userID) {
		return nil, ErrDeliveryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load delivery: %w", err)
	}
	return delivery, nil
}

func (s *Service) deliveryByRequest(ctx context.Context, userID, requestID uuid.UUID) (*db.PaperLetterDelivery, error) {
	return s.repository.GetPaperLetterDeliveryByUserAndClientRequestID(ctx, db.GetPaperLetterDeliveryByUserAndClientRequestIDParams{UserID: userID, ClientRequestID: requestID})
}

func (p *DeliveryProcessor) ProcessDelivery(ctx context.Context, deliveryID uuid.UUID) error {
	delivery, err := p.repository.GetPaperLetterDeliveryByID(ctx, db.GetPaperLetterDeliveryByIDParams{ID: deliveryID})
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrDeliveryNotFound
	}
	if err != nil {
		return fmt.Errorf("load delivery: %w", err)
	}
	if delivery.Status == db.PaperLetterDeliveryStatusSent || delivery.Status == db.PaperLetterDeliveryStatusFailed || delivery.Status == db.PaperLetterDeliveryStatusManualReview {
		return nil
	}
	delivery, err = p.repository.MarkPaperLetterDeliveryProcessing(ctx, db.MarkPaperLetterDeliveryProcessingParams{ID: deliveryID})
	if err != nil {
		return fmt.Errorf("mark delivery processing: %w", err)
	}
	now := p.now()
	if !now.Before(delivery.ProviderSubmissionStartedAt.Add(24 * time.Hour)) {
		cutoff := now.Add(-24 * time.Hour)
		_, err := p.repository.MarkPaperLetterDeliveryManualReview(ctx, db.MarkPaperLetterDeliveryManualReviewParams{
			ID: delivery.ID, Cutoff: &cutoff,
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("mark delivery for manual review: %w", err)
		}
		return nil
	}
	upload, err := p.repository.GetPaperLetterUploadByIDAndUserID(ctx, db.GetPaperLetterUploadByIDAndUserIDParams{ID: delivery.UploadID, UserID: delivery.UserID})
	if err != nil {
		return fmt.Errorf("load delivery upload: %w", err)
	}
	data, err := p.storage.DownloadObject(ctx, delivery.UserID, paperLetterObjectKey(upload.ID), p.maxPDFBytes)
	if err != nil {
		return fmt.Errorf("download delivery PDF: %w", err)
	}
	sum := sha256.Sum256(data)
	pageCount, pdfErr := validateUploadedPDF(data)
	if int64(len(data)) != upload.SizeBytes || hex.EncodeToString(sum[:]) != upload.Sha256 || pdfErr != nil || int32(pageCount) != upload.PageCount || upload.PageCount != delivery.PageCount {
		return p.failDelivery(ctx, delivery, "stored_pdf_mismatch")
	}
	pingenID, err := p.provider.Submit(ctx, data, upload.SanitizedFilename, delivery.ID.String(), delivery.DeliveryProduct, delivery.AddressPosition, delivery.PrintMode, delivery.PrintSpectrum)
	if err != nil {
		var temporary interface{ Temporary() bool }
		if !errors.As(err, &temporary) || !temporary.Temporary() {
			message := "pingen_terminal_failure"
			var safe interface{ SafeMessage() string }
			if errors.As(err, &safe) {
				message = safe.SafeMessage()
			}
			return p.failDelivery(ctx, delivery, message)
		}
		return fmt.Errorf("submit delivery: %w", err)
	}
	if _, err := p.repository.MarkPaperLetterDeliverySent(ctx, db.MarkPaperLetterDeliverySentParams{ID: delivery.ID, PingenID: pgtype.Text{String: pingenID, Valid: true}}); err != nil {
		return fmt.Errorf("mark delivery sent: %w", err)
	}
	return nil
}

func IsPermanentDelivery(err error) bool {
	return errors.Is(err, ErrDeliveryNotFound)
}

func (p *DeliveryProcessor) failDelivery(ctx context.Context, delivery *db.PaperLetterDelivery, message string) error {
	tx, err := p.repository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin delivery failure transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if _, err := p.repository.WithTx(tx).FailPaperLetterDelivery(ctx, db.FailPaperLetterDeliveryParams{
		ID: delivery.ID, ErrorMessage: pgtype.Text{String: message, Valid: true},
	}); errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if err != nil {
		return fmt.Errorf("mark delivery failed: %w", err)
	}
	if _, err := p.refunder.RefundTx(ctx, tx, users_service.RefundTxInput{
		UserID: delivery.UserID, Amount: delivery.BpCost, Description: "Papierbrief delivery refund",
		RelatedEntityID: delivery.ID, RelatedEntityType: db.RelatedEntityTypePaperLetterDelivery,
	}); err != nil {
		return fmt.Errorf("refund delivery BuddyPoints: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit delivery failure: %w", err)
	}
	return nil
}

func isDeliveryRequestConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == deliveryRequestConstraint
}
