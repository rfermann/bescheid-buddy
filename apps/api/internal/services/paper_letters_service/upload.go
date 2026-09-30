package paper_letters_service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bored-bots/bored-template/internal/services/document_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

var (
	ErrPaperLetterUploadNotFound = errors.New("paper letter upload was not found")
	ErrUploadedPDFNotFound       = errors.New("uploaded PDF was not found in storage")
	ErrInvalidUploadFilename     = errors.New("invalid paper letter upload filename")
	ErrPDFTooLarge               = errors.New("paper letter PDF exceeds the maximum allowed size")
	ErrInvalidPDF                = errors.New("invalid paper letter PDF")
)

type PaperLetterStorage interface {
	PresignObjectUpload(ctx context.Context, userID uuid.UUID, objectKey string) (string, error)
	PresignObjectDownload(ctx context.Context, userID uuid.UUID, objectKey, filename string, attachment bool) (string, error)
	StatObjectSize(ctx context.Context, userID uuid.UUID, objectKey string) (int64, error)
	DownloadObject(ctx context.Context, userID uuid.UUID, objectKey string, maxBytes int64) ([]byte, error)
}

type CreateUploadResult struct {
	UploadID  uuid.UUID
	UploadURL string
}

type ConfirmUploadResult struct {
	Upload        *db.PaperLetterUpload
	DeliveryPrice int
}

func paperLetterObjectKey(uploadID uuid.UUID) string {
	return "paper-letters/" + uploadID.String() + ".pdf"
}

func (s *Service) CreateUpload(ctx context.Context, userID uuid.UUID, filename string) (*CreateUploadResult, error) {
	filename, err := sanitizeUploadFilename(filename)
	if err != nil {
		return nil, err
	}
	uploadID := uuid.New()
	_, err = s.repository.CreatePaperLetterUpload(ctx, db.CreatePaperLetterUploadParams{
		ID: uploadID, UserID: userID, SanitizedFilename: filename,
	})
	if err != nil {
		return nil, fmt.Errorf("create paper letter upload: %w", err)
	}
	uploadURL, err := s.storage.PresignObjectUpload(ctx, userID, paperLetterObjectKey(uploadID))
	if err != nil {
		return nil, fmt.Errorf("presign paper letter upload: %w", err)
	}
	return &CreateUploadResult{UploadID: uploadID, UploadURL: uploadURL}, nil
}

func (s *Service) ConfirmUpload(ctx context.Context, userID, uploadID uuid.UUID) (*ConfirmUploadResult, error) {
	upload, err := s.repository.GetPaperLetterUploadByIDAndUserID(ctx, db.GetPaperLetterUploadByIDAndUserIDParams{
		ID: uploadID, UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPaperLetterUploadNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load paper letter upload: %w", err)
	}
	if upload.Status == db.PaperLetterUploadStatusConfirmed {
		return confirmedUploadResult(upload), nil
	}

	objectKey := paperLetterObjectKey(uploadID)
	size, err := s.storage.StatObjectSize(ctx, userID, objectKey)
	if errors.Is(err, document_service.ErrObjectNotFound) {
		return nil, ErrUploadedPDFNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("stat paper letter upload: %w", err)
	}
	if size > s.maxPDFBytes {
		return nil, ErrPDFTooLarge
	}
	data, err := s.storage.DownloadObject(ctx, userID, objectKey, s.maxPDFBytes)
	if errors.Is(err, document_service.ErrObjectNotFound) {
		return nil, ErrUploadedPDFNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("download paper letter upload: %w", err)
	}
	if int64(len(data)) != size || int64(len(data)) > s.maxPDFBytes {
		return nil, ErrPDFTooLarge
	}
	pageCount, err := validateUploadedPDF(data)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(data)
	upload, err = s.repository.ConfirmPaperLetterUpload(ctx, db.ConfirmPaperLetterUploadParams{
		ID: uploadID, UserID: userID, SizeBytes: size, PageCount: int32(pageCount), Sha256: hex.EncodeToString(sum[:]),
	})
	if err != nil {
		return nil, fmt.Errorf("confirm paper letter upload: %w", err)
	}
	return confirmedUploadResult(upload), nil
}

func confirmedUploadResult(upload *db.PaperLetterUpload) *ConfirmUploadResult {
	return &ConfirmUploadResult{Upload: upload, DeliveryPrice: deliveryPrice(int(upload.PageCount))}
}

func deliveryPrice(pageCount int) int {
	if pageCount < 1 {
		return 0
	}
	return 50 + (pageCount-1)*20
}

func validateUploadedPDF(data []byte) (int, error) {
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return 0, ErrInvalidPDF
	}
	configuration := model.NewDefaultConfiguration()
	configuration.Offline = true
	if err := api.Validate(bytes.NewReader(data), configuration); err != nil {
		return 0, ErrInvalidPDF
	}
	pageCount, err := api.PageCount(bytes.NewReader(data), configuration)
	if err != nil || pageCount < 1 {
		return 0, ErrInvalidPDF
	}
	return pageCount, nil
}

func sanitizeUploadFilename(value string) (string, error) {
	value = path.Base(strings.ReplaceAll(strings.TrimSpace(value), "\\", "/"))
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return '_'
		}
		return r
	}, value)
	if value == "." || value == "" || !strings.EqualFold(path.Ext(value), ".pdf") || utf8.RuneCountInString(value) > 255 {
		return "", ErrInvalidUploadFilename
	}
	return value, nil
}
