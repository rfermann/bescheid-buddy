package paper_letters_service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/services/document_service"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_pdf_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const defaultPaperLetterMaxPDFBytes = int64(15_728_640)

type uploadRepositoryFake struct {
	cases_repository.Repository
	uploads map[uuid.UUID]*db.PaperLetterUpload
}

func newUploadRepositoryFake() *uploadRepositoryFake {
	return &uploadRepositoryFake{uploads: make(map[uuid.UUID]*db.PaperLetterUpload)}
}

func (r *uploadRepositoryFake) CreatePaperLetterUpload(_ context.Context, input db.CreatePaperLetterUploadParams) (*db.PaperLetterUpload, error) {
	now := time.Now()
	upload := &db.PaperLetterUpload{
		ID: input.ID, UserID: input.UserID, SanitizedFilename: input.SanitizedFilename,
		Status: db.PaperLetterUploadStatusPending, CreatedAt: now, UpdatedAt: now,
	}
	r.uploads[upload.ID] = upload
	return upload, nil
}

func (r *uploadRepositoryFake) GetPaperLetterUploadByIDAndUserID(_ context.Context, input db.GetPaperLetterUploadByIDAndUserIDParams) (*db.PaperLetterUpload, error) {
	upload := r.uploads[input.ID]
	if upload == nil || upload.UserID != input.UserID {
		return nil, pgx.ErrNoRows
	}
	copy := *upload
	return &copy, nil
}

func (r *uploadRepositoryFake) ConfirmPaperLetterUpload(_ context.Context, input db.ConfirmPaperLetterUploadParams) (*db.PaperLetterUpload, error) {
	upload := r.uploads[input.ID]
	if upload == nil || upload.UserID != input.UserID || upload.Status != db.PaperLetterUploadStatusPending {
		return nil, pgx.ErrNoRows
	}
	now := time.Now()
	upload.Status = db.PaperLetterUploadStatusConfirmed
	upload.SizeBytes = input.SizeBytes
	upload.PageCount = input.PageCount
	upload.Sha256 = input.Sha256
	upload.ConfirmedAt = &now
	upload.UpdatedAt = now
	copy := *upload
	return &copy, nil
}

type uploadStorageFake struct {
	url         string
	data        []byte
	statErr     error
	downloadErr error
	userID      uuid.UUID
	objectKey   string
	maxBytes    int64
	statCalls   int
}

func (s *uploadStorageFake) PresignObjectUpload(_ context.Context, userID uuid.UUID, objectKey string) (string, error) {
	s.userID, s.objectKey = userID, objectKey
	return s.url, nil
}

func (s *uploadStorageFake) StatObjectSize(_ context.Context, userID uuid.UUID, objectKey string) (int64, error) {
	s.userID, s.objectKey = userID, objectKey
	s.statCalls++
	return int64(len(s.data)), s.statErr
}

func (s *uploadStorageFake) DownloadObject(_ context.Context, userID uuid.UUID, objectKey string, maxBytes int64) ([]byte, error) {
	s.userID, s.objectKey, s.maxBytes = userID, objectKey, maxBytes
	return s.data, s.downloadErr
}

func uploadTestService(repository *uploadRepositoryFake, storage *uploadStorageFake, maxBytes int64) *Service {
	return &Service{repository: repository, storage: storage, maxPDFBytes: maxBytes}
}

func generatedUploadPDF(t *testing.T) []byte {
	t.Helper()
	result, err := paper_letters_pdf_service.New().Generate(paper_letters_pdf_service.GenerateInput{
		GermanText: "Hiermit widerspreche ich dem Bescheid.",
		Subject:    "Widerspruch",
		Sender:     paper_letters_pdf_service.Party{Name: "René Müller", Address: "Teststraße 1\n10115 Berlin"},
		Recipient:  paper_letters_pdf_service.Party{Name: "Behörde", Address: "Amtstraße 2\n20095 Hamburg"},
		Filename:   "brief.pdf",
		Now:        time.Date(2026, time.September, 15, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	return result.Bytes
}

func TestCreateUploadSanitizesFilenameAndOwnsObjectKey(t *testing.T) {
	repository := newUploadRepositoryFake()
	storage := &uploadStorageFake{url: "https://storage.example/upload"}
	service := uploadTestService(repository, storage, defaultPaperLetterMaxPDFBytes)
	userID := uuid.New()

	result, err := service.CreateUpload(context.Background(), userID, "../letter\r.pdf")

	require.NoError(t, err)
	assert.Equal(t, "letter_.pdf", repository.uploads[result.UploadID].SanitizedFilename)
	assert.Equal(t, userID, storage.userID)
	assert.Equal(t, "paper-letters/"+result.UploadID.String()+".pdf", storage.objectKey)
	assert.Equal(t, "https://storage.example/upload", result.UploadURL)

	_, err = service.CreateUpload(context.Background(), userID, "letter.txt")
	assert.ErrorIs(t, err, ErrInvalidUploadFilename)
}

func TestConfirmUploadValidatesAndPersistsGeneratedPDF(t *testing.T) {
	data := generatedUploadPDF(t)
	repository := newUploadRepositoryFake()
	storage := &uploadStorageFake{data: data}
	service := uploadTestService(repository, storage, int64(len(data)))
	userID, uploadID := uuid.New(), uuid.New()
	repository.uploads[uploadID] = &db.PaperLetterUpload{
		ID: uploadID, UserID: userID, SanitizedFilename: "letter.pdf", Status: db.PaperLetterUploadStatusPending,
	}

	result, err := service.ConfirmUpload(context.Background(), userID, uploadID)

	require.NoError(t, err)
	assert.Equal(t, db.PaperLetterUploadStatusConfirmed, result.Upload.Status)
	assert.Equal(t, int64(len(data)), result.Upload.SizeBytes)
	assert.Equal(t, int32(1), result.Upload.PageCount)
	sum := sha256.Sum256(data)
	assert.Equal(t, hex.EncodeToString(sum[:]), result.Upload.Sha256)
	assert.Equal(t, 50, result.DeliveryPrice)
	assert.Equal(t, int64(len(data)), storage.maxBytes)

	_, err = service.ConfirmUpload(context.Background(), userID, uploadID)
	require.NoError(t, err)
	assert.Equal(t, 1, storage.statCalls, "confirmed uploads must be idempotent")
}

func TestConfirmUploadRejectsWrongOwnerMissingOversizeAndInvalidPDF(t *testing.T) {
	data := generatedUploadPDF(t)
	userID, uploadID := uuid.New(), uuid.New()
	newPending := func() *uploadRepositoryFake {
		repository := newUploadRepositoryFake()
		repository.uploads[uploadID] = &db.PaperLetterUpload{
			ID: uploadID, UserID: userID, SanitizedFilename: "letter.pdf", Status: db.PaperLetterUploadStatusPending,
		}
		return repository
	}

	t.Run("ownership", func(t *testing.T) {
		storage := &uploadStorageFake{data: data}
		_, err := uploadTestService(newPending(), storage, int64(len(data))).ConfirmUpload(context.Background(), uuid.New(), uploadID)
		assert.ErrorIs(t, err, ErrPaperLetterUploadNotFound)
		assert.Zero(t, storage.statCalls)
	})

	t.Run("missing object", func(t *testing.T) {
		storage := &uploadStorageFake{statErr: document_service.ErrObjectNotFound}
		_, err := uploadTestService(newPending(), storage, defaultPaperLetterMaxPDFBytes).ConfirmUpload(context.Background(), userID, uploadID)
		assert.ErrorIs(t, err, ErrUploadedPDFNotFound)
	})

	t.Run("oversize", func(t *testing.T) {
		storage := &uploadStorageFake{data: data}
		_, err := uploadTestService(newPending(), storage, int64(len(data)-1)).ConfirmUpload(context.Background(), userID, uploadID)
		assert.ErrorIs(t, err, ErrPDFTooLarge)
	})

	t.Run("invalid PDF", func(t *testing.T) {
		storage := &uploadStorageFake{data: []byte("%PDF-not-really-a-pdf")}
		_, err := uploadTestService(newPending(), storage, defaultPaperLetterMaxPDFBytes).ConfirmUpload(context.Background(), userID, uploadID)
		assert.ErrorIs(t, err, ErrInvalidPDF)
	})
}

func TestConfirmUploadAllowsExactConfiguredBoundary(t *testing.T) {
	data := generatedUploadPDF(t)
	require.Less(t, int64(len(data)), defaultPaperLetterMaxPDFBytes)
	data = append(data, bytes.Repeat([]byte("\n"), int(defaultPaperLetterMaxPDFBytes)-len(data))...)
	repository := newUploadRepositoryFake()
	storage := &uploadStorageFake{data: data}
	service := uploadTestService(repository, storage, defaultPaperLetterMaxPDFBytes)
	userID, uploadID := uuid.New(), uuid.New()
	repository.uploads[uploadID] = &db.PaperLetterUpload{
		ID: uploadID, UserID: userID, SanitizedFilename: "letter.pdf", Status: db.PaperLetterUploadStatusPending,
	}

	result, err := service.ConfirmUpload(context.Background(), userID, uploadID)

	require.NoError(t, err)
	assert.Equal(t, defaultPaperLetterMaxPDFBytes, result.Upload.SizeBytes)
}

func TestDeliveryPrice(t *testing.T) {
	assert.Equal(t, 0, deliveryPrice(0))
	assert.Equal(t, 50, deliveryPrice(1))
	assert.Equal(t, 90, deliveryPrice(3))
}
