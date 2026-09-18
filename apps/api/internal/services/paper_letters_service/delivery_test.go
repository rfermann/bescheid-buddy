package paper_letters_service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/services/users_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deliveryTxFake struct {
	pgx.Tx
	mu                    sync.Mutex
	committed, rolledBack bool
}

func (tx *deliveryTxFake) Commit(context.Context) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tx.committed = true
	return nil
}
func (tx *deliveryTxFake) Rollback(context.Context) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if !tx.committed {
		tx.rolledBack = true
	}
	return nil
}

type deliveryRepositoryFake struct {
	cases_repository.Repository
	tx              *deliveryTxFake
	upload          *db.PaperLetterUpload
	deliveries      map[uuid.UUID]*db.PaperLetterDelivery
	createErr       error
	skipFirstLookup bool
	lookupCalls     int
	mu              sync.Mutex
}

func newDeliveryRepositoryFake(upload *db.PaperLetterUpload) *deliveryRepositoryFake {
	return &deliveryRepositoryFake{tx: &deliveryTxFake{}, upload: upload, deliveries: map[uuid.UUID]*db.PaperLetterDelivery{}}
}

func (r *deliveryRepositoryFake) BeginTx(context.Context) (pgx.Tx, error)   { return r.tx, nil }
func (r *deliveryRepositoryFake) WithTx(pgx.Tx) cases_repository.Repository { return r }
func (r *deliveryRepositoryFake) GetPaperLetterUploadByIDAndUserID(_ context.Context, p db.GetPaperLetterUploadByIDAndUserIDParams) (*db.PaperLetterUpload, error) {
	if r.upload == nil || r.upload.ID != p.ID || r.upload.UserID != p.UserID {
		return nil, pgx.ErrNoRows
	}
	return r.upload, nil
}
func (r *deliveryRepositoryFake) CreatePaperLetterDelivery(_ context.Context, p db.CreatePaperLetterDeliveryParams) (*db.PaperLetterDelivery, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	delivery := &db.PaperLetterDelivery{
		ID: p.ID, UserID: p.UserID, UploadID: p.UploadID, ClientRequestID: p.ClientRequestID,
		Status: db.PaperLetterDeliveryStatusPending, PageCount: p.PageCount, BpCost: p.BpCost,
		DeliveryProduct: p.DeliveryProduct, AddressPosition: p.AddressPosition,
		PrintMode: p.PrintMode, PrintSpectrum: p.PrintSpectrum,
	}
	r.deliveries[delivery.ID] = delivery
	return delivery, nil
}
func (r *deliveryRepositoryFake) GetPaperLetterDeliveryByID(_ context.Context, p db.GetPaperLetterDeliveryByIDParams) (*db.PaperLetterDelivery, error) {
	delivery := r.deliveries[p.ID]
	if delivery == nil {
		return nil, pgx.ErrNoRows
	}
	return delivery, nil
}
func (r *deliveryRepositoryFake) GetPaperLetterDeliveryByUserAndClientRequestID(_ context.Context, p db.GetPaperLetterDeliveryByUserAndClientRequestIDParams) (*db.PaperLetterDelivery, error) {
	r.lookupCalls++
	if r.skipFirstLookup && r.lookupCalls == 1 {
		return nil, pgx.ErrNoRows
	}
	for _, delivery := range r.deliveries {
		if delivery.UserID == p.UserID && delivery.ClientRequestID == p.ClientRequestID {
			return delivery, nil
		}
	}
	return nil, pgx.ErrNoRows
}
func (r *deliveryRepositoryFake) SetPaperLetterDeliveryJobID(_ context.Context, p db.SetPaperLetterDeliveryJobIDParams) error {
	r.deliveries[p.ID].JobID = p.JobID
	return nil
}
func (r *deliveryRepositoryFake) MarkPaperLetterDeliveryProcessing(_ context.Context, p db.MarkPaperLetterDeliveryProcessingParams) (*db.PaperLetterDelivery, error) {
	delivery := r.deliveries[p.ID]
	delivery.Status = db.PaperLetterDeliveryStatusProcessing
	if delivery.ProviderSubmissionStartedAt == nil {
		now := time.Now()
		delivery.ProviderSubmissionStartedAt = &now
	}
	return delivery, nil
}
func (r *deliveryRepositoryFake) MarkPaperLetterDeliverySent(_ context.Context, p db.MarkPaperLetterDeliverySentParams) (*db.PaperLetterDelivery, error) {
	delivery := r.deliveries[p.ID]
	delivery.Status, delivery.PingenID = db.PaperLetterDeliveryStatusSent, p.PingenID
	return delivery, nil
}
func (r *deliveryRepositoryFake) MarkPaperLetterDeliveryManualReview(_ context.Context, p db.MarkPaperLetterDeliveryManualReviewParams) (*db.PaperLetterDelivery, error) {
	delivery := r.deliveries[p.ID]
	if delivery.Status != db.PaperLetterDeliveryStatusProcessing || delivery.ProviderSubmissionStartedAt.After(*p.Cutoff) {
		return nil, pgx.ErrNoRows
	}
	delivery.Status = db.PaperLetterDeliveryStatusManualReview
	return delivery, nil
}
func (r *deliveryRepositoryFake) FailPaperLetterDelivery(_ context.Context, p db.FailPaperLetterDeliveryParams) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delivery := r.deliveries[p.ID]
	if delivery.Status != db.PaperLetterDeliveryStatusProcessing {
		return uuid.Nil, pgx.ErrNoRows
	}
	delivery.Status, delivery.ErrorMessage = db.PaperLetterDeliveryStatusFailed, p.ErrorMessage
	return delivery.ID, nil
}

type deliveryPointsFake struct {
	input       users_service.SpendTxInput
	refundInput users_service.RefundTxInput
	refunds     int
	err         error
	mu          sync.Mutex
}

func (p *deliveryPointsFake) RefundTx(_ context.Context, _ pgx.Tx, input users_service.RefundTxInput) (*users_service.RefundTxResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.refundInput, p.refunds = input, p.refunds+1
	return &users_service.RefundTxResult{}, p.err
}

func (p *deliveryPointsFake) SpendTx(_ context.Context, _ pgx.Tx, input users_service.SpendTxInput) (*users_service.SpendTxResult, error) {
	p.input = input
	return &users_service.SpendTxResult{}, p.err
}

type deliveryQueueFake struct{ err error }

func (q deliveryQueueFake) QueueDocumentReplyGenerationJobTx(context.Context, pgx.Tx, uuid.UUID) (int64, error) {
	return 0, nil
}
func (q deliveryQueueFake) QueuePaperLetterDeliveryJobTx(context.Context, pgx.Tx, uuid.UUID) (int64, error) {
	return 42, q.err
}

func TestCreateDelivery(t *testing.T) {
	userID, uploadID, requestID := uuid.New(), uuid.New(), uuid.New()
	confirmed := &db.PaperLetterUpload{ID: uploadID, UserID: userID, Status: db.PaperLetterUploadStatusConfirmed, PageCount: 3}
	newService := func(repo *deliveryRepositoryFake, points *deliveryPointsFake, queue deliveryQueueFake) *Service {
		return &Service{repository: repo, points: points, queue: queue, pingen: config.PingenConfig{
			Product: "fast", AddressPosition: "left", PrintMode: "simplex", PrintSpectrum: "color",
		}}
	}

	t.Run("charges creates and enqueues once", func(t *testing.T) {
		repo, points := newDeliveryRepositoryFake(confirmed), &deliveryPointsFake{}
		service := newService(repo, points, deliveryQueueFake{})
		result, err := service.CreateDelivery(context.Background(), CreateDeliveryInput{UserID: userID, UploadID: uploadID, ClientRequestID: requestID})
		require.NoError(t, err)
		assert.Equal(t, int32(90), result.Delivery.BpCost)
		assert.Equal(t, int64(42), result.Delivery.JobID.Int64)
		assert.Equal(t, db.RelatedEntityTypePaperLetterDelivery, points.input.RelatedEntityType)
		assert.Equal(t, int32(90), points.input.Amount)
		assert.True(t, repo.tx.committed)
		_, err = service.GetDelivery(context.Background(), uuid.New(), result.Delivery.ID)
		assert.ErrorIs(t, err, ErrDeliveryNotFound)

		again, err := service.CreateDelivery(context.Background(), CreateDeliveryInput{UserID: userID, UploadID: uploadID, ClientRequestID: requestID})
		require.NoError(t, err)
		assert.True(t, again.Existing)
	})

	tests := []struct {
		name   string
		upload *db.PaperLetterUpload
		points error
		queue  error
		want   error
	}{
		{"ownership", &db.PaperLetterUpload{ID: uploadID, UserID: uuid.New()}, nil, nil, ErrPaperLetterUploadNotFound},
		{"unconfirmed", &db.PaperLetterUpload{ID: uploadID, UserID: userID, Status: db.PaperLetterUploadStatusPending}, nil, nil, ErrUploadNotConfirmed},
		{"insufficient points", confirmed, users_service.ErrInsufficientBuddyPoints, nil, ErrInsufficientBuddyPoints},
		{"queue rollback", confirmed, nil, errors.New("queue failed"), errors.New("enqueue delivery")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newDeliveryRepositoryFake(test.upload)
			_, err := newService(repo, &deliveryPointsFake{err: test.points}, deliveryQueueFake{err: test.queue}).CreateDelivery(
				context.Background(), CreateDeliveryInput{UserID: userID, UploadID: uploadID, ClientRequestID: uuid.New()},
			)
			if test.name == "queue rollback" {
				assert.ErrorContains(t, err, test.want.Error())
			} else {
				assert.ErrorIs(t, err, test.want)
			}
			assert.True(t, repo.tx.rolledBack)
		})
	}

	t.Run("concurrent duplicate rolls back charge", func(t *testing.T) {
		repo, points := newDeliveryRepositoryFake(confirmed), &deliveryPointsFake{}
		existing := &db.PaperLetterDelivery{ID: uuid.New(), UserID: userID, ClientRequestID: requestID}
		repo.deliveries[existing.ID] = existing
		repo.createErr = &pgconn.PgError{Code: "23505", ConstraintName: deliveryRequestConstraint}
		repo.skipFirstLookup = true
		result, err := newService(repo, points, deliveryQueueFake{}).CreateDelivery(context.Background(), CreateDeliveryInput{
			UserID: userID, UploadID: uploadID, ClientRequestID: requestID,
		})
		require.NoError(t, err)
		assert.True(t, result.Existing)
	})
}

type deliveryProviderFake struct {
	idempotencyKey string
	called         bool
	err            error
}

func (p *deliveryProviderFake) Submit(_ context.Context, _ []byte, _, idempotencyKey, _, _, _, _ string) (string, error) {
	p.called, p.idempotencyKey = true, idempotencyKey
	return "pingen-id", p.err
}

type temporaryDeliveryError struct{}

func (temporaryDeliveryError) Error() string   { return "temporary" }
func (temporaryDeliveryError) Temporary() bool { return true }

type terminalDeliveryError struct{}

func (terminalDeliveryError) Error() string       { return "recipient secret" }
func (terminalDeliveryError) Temporary() bool     { return false }
func (terminalDeliveryError) SafeMessage() string { return "pingen_status_422_address_invalid" }

func TestProcessDeliveryVerifiesPDFAndMarksSent(t *testing.T) {
	data := generatedUploadPDF(t)
	sum := sha256.Sum256(data)
	userID, uploadID, deliveryID := uuid.New(), uuid.New(), uuid.New()
	upload := &db.PaperLetterUpload{
		ID: uploadID, UserID: userID, SanitizedFilename: "letter.pdf", Status: db.PaperLetterUploadStatusConfirmed,
		SizeBytes: int64(len(data)), PageCount: 1, Sha256: hex.EncodeToString(sum[:]),
	}
	repo := newDeliveryRepositoryFake(upload)
	repo.deliveries[deliveryID] = &db.PaperLetterDelivery{
		ID: deliveryID, UserID: userID, UploadID: uploadID, Status: db.PaperLetterDeliveryStatusPending,
		PageCount: 1, DeliveryProduct: "fast", AddressPosition: "left", PrintMode: "simplex", PrintSpectrum: "color",
	}
	provider := &deliveryProviderFake{}
	points := &deliveryPointsFake{}
	processor := NewDeliveryProcessor(repo, &uploadStorageFake{data: data}, provider, points, int64(len(data)))
	require.NoError(t, processor.ProcessDelivery(context.Background(), deliveryID))
	assert.Equal(t, db.PaperLetterDeliveryStatusSent, repo.deliveries[deliveryID].Status)
	assert.Equal(t, "pingen-id", repo.deliveries[deliveryID].PingenID.String)
	assert.Equal(t, deliveryID.String(), provider.idempotencyKey)
	assert.Zero(t, points.refunds)
	provider.err = terminalDeliveryError{}
	require.NoError(t, processor.ProcessDelivery(context.Background(), deliveryID))
	assert.Equal(t, db.PaperLetterDeliveryStatusSent, repo.deliveries[deliveryID].Status)
	assert.Zero(t, points.refunds)

	now := time.Now()
	processor.now = func() time.Time { return now }
	delivery := repo.deliveries[deliveryID]
	delivery.Status = db.PaperLetterDeliveryStatusProcessing
	delivery.ProviderSubmissionStartedAt = ptr(now.Add(-24*time.Hour + time.Nanosecond))
	provider.err = temporaryDeliveryError{}
	err := processor.ProcessDelivery(context.Background(), deliveryID)
	assert.Error(t, err)
	assert.False(t, IsPermanentDelivery(err))
	assert.Equal(t, deliveryID.String(), provider.idempotencyKey)

	delivery.ProviderSubmissionStartedAt = ptr(now.Add(-24 * time.Hour))
	provider.called = false
	err = processor.ProcessDelivery(context.Background(), deliveryID)
	require.NoError(t, err)
	assert.Equal(t, db.PaperLetterDeliveryStatusManualReview, delivery.Status)
	assert.False(t, provider.called)
	require.NoError(t, processor.ProcessDelivery(context.Background(), deliveryID))
	assert.False(t, provider.called)

	delivery.Status = db.PaperLetterDeliveryStatusPending
	delivery.ProviderSubmissionStartedAt = ptr(now)
	upload.Sha256 = "bad"
	err = processor.ProcessDelivery(context.Background(), deliveryID)
	require.NoError(t, err)
	assert.Equal(t, db.PaperLetterDeliveryStatusFailed, delivery.Status)
	assert.Equal(t, 1, points.refunds)
}

func TestTerminalDeliveryRefundsOnce(t *testing.T) {
	data := generatedUploadPDF(t)
	sum := sha256.Sum256(data)
	userID, uploadID, deliveryID := uuid.New(), uuid.New(), uuid.New()
	upload := &db.PaperLetterUpload{ID: uploadID, UserID: userID, SanitizedFilename: "letter.pdf", SizeBytes: int64(len(data)), PageCount: 1, Sha256: hex.EncodeToString(sum[:])}
	delivery := &db.PaperLetterDelivery{ID: deliveryID, UserID: userID, UploadID: uploadID, Status: db.PaperLetterDeliveryStatusPending, PageCount: 1, BpCost: 50}
	repo, points := newDeliveryRepositoryFake(upload), &deliveryPointsFake{}
	repo.deliveries[deliveryID] = delivery
	processor := NewDeliveryProcessor(repo, &uploadStorageFake{data: data}, &deliveryProviderFake{err: terminalDeliveryError{}}, points, int64(len(data)))

	require.NoError(t, processor.ProcessDelivery(context.Background(), deliveryID))
	assert.Equal(t, db.PaperLetterDeliveryStatusFailed, delivery.Status)
	assert.Equal(t, "pingen_status_422_address_invalid", delivery.ErrorMessage.String)
	assert.NotContains(t, delivery.ErrorMessage.String, "recipient secret")
	assert.Equal(t, users_service.RefundTxInput{UserID: userID, Amount: 50, Description: "Papierbrief delivery refund", RelatedEntityID: deliveryID, RelatedEntityType: db.RelatedEntityTypePaperLetterDelivery}, points.refundInput)

	delivery.Status = db.PaperLetterDeliveryStatusProcessing
	points.refunds = 0
	errs := make(chan error, 2)
	for range 2 {
		go func() {
			errs <- processor.failDelivery(context.Background(), delivery, "pingen_status_422_address_invalid")
		}()
	}
	require.NoError(t, <-errs)
	require.NoError(t, <-errs)
	assert.Equal(t, 1, points.refunds)
}

func ptr[T any](value T) *T { return &value }
