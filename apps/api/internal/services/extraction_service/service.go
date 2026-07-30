package extraction_service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/services/document_service"
	"github.com/bored-bots/bored-template/internal/services/vision_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

const visionPDFChunkPages = 5

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrFileTooLarge        = errors.New("file exceeds the maximum allowed size")
	ErrNoTextExtracted     = errors.New("no text could be extracted from the document")
	ErrEncryptedPDF        = errors.New("pdf is password-protected")
	ErrCorruptPDF          = errors.New("pdf could not be parsed")
)

type Service struct {
	source     *document_service.Service
	recognizer *vision_service.Service
	repository cases_repository.Repository
	maxBytes   int64
}

func New(source *document_service.Service, recognizer *vision_service.Service, repository cases_repository.Repository, maxBytes int64) *Service {
	return &Service{
		source:     source,
		recognizer: recognizer,
		repository: repository,
		maxBytes:   maxBytes,
	}
}

func (s *Service) ExtractDocument(ctx context.Context, userID, caseID, documentID uuid.UUID) error {
	if _, err := s.repository.MarkDocumentExtractionProcessing(ctx, db.MarkDocumentExtractionProcessingParams{
		DocumentID: documentID,
	}); err != nil {
		return err
	}

	text, provider, mime, pageCount, err := s.extract(ctx, userID, caseID, documentID)
	if err != nil {
		if IsPermanent(err) {
			return s.fail(ctx, documentID, err)
		}
		return err
	}

	return s.complete(ctx, documentID, text, provider, mime, pageCount)
}

func (s *Service) extract(ctx context.Context, userID, caseID, documentID uuid.UUID) (text, provider, mime string, pageCount int, err error) {
	data, err := s.source.DownloadDocument(ctx, userID, caseID, documentID, s.maxBytes)
	if err != nil {
		return "", "", "", 0, err
	}
	if int64(len(data)) > s.maxBytes {
		return "", "", "", 0, ErrFileTooLarge
	}

	mime = sniffMIME(data)
	switch mime {
	case "text/plain":
		return normalize(string(data)), "plaintext", mime, 0, nil

	case "image/jpeg", "image/png", "image/webp":
		raw, rerr := s.recognizer.DetectImageText(ctx, data)
		if rerr != nil {
			return "", "", "", 0, rerr
		}
		text = normalize(raw)
		if text == "" {
			return "", "", "", 0, ErrNoTextExtracted
		}
		return text, "google_vision", mime, 0, nil

	case "application/pdf":
		raw, pages, perr := s.extractPDF(ctx, data)
		if perr != nil {
			return "", "", "", 0, perr
		}
		text = normalize(raw)
		if text == "" {
			return "", "", "", 0, ErrNoTextExtracted
		}
		return text, "google_vision", mime, pages, nil

	default:
		return "", "", "", 0, fmt.Errorf("%w: %s", ErrUnsupportedFileType, mime)
	}
}

func (s *Service) extractPDF(ctx context.Context, data []byte) (string, int, error) {
	conf := model.NewDefaultConfiguration()
	conf.Offline = true

	if err := api.Validate(bytes.NewReader(data), conf); err != nil {
		if isEncryptedPDFErr(err) {
			return "", 0, ErrEncryptedPDF
		}
		return "", 0, fmt.Errorf("%w: %v", ErrCorruptPDF, err)
	}

	pageCount, err := api.PageCount(bytes.NewReader(data), conf)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", ErrCorruptPDF, err)
	}
	if pageCount == 0 {
		return "", 0, ErrNoTextExtracted
	}

	var b strings.Builder
	for _, pages := range chunkPages(pageCount, visionPDFChunkPages) {
		chunkText, cerr := s.recognizer.DetectPDFText(ctx, data, pages)
		if cerr != nil {
			return "", pageCount, cerr
		}
		if chunkText == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(chunkText)
	}

	return b.String(), pageCount, nil
}

func chunkPages(pageCount, chunkSize int) [][]int {
	var chunks [][]int
	for start := 1; start <= pageCount; start += chunkSize {
		end := start + chunkSize - 1
		if end > pageCount {
			end = pageCount
		}
		pages := make([]int, 0, end-start+1)
		for p := start; p <= end; p++ {
			pages = append(pages, p)
		}
		chunks = append(chunks, pages)
	}
	return chunks
}

func (s *Service) complete(ctx context.Context, documentID uuid.UUID, text, provider, mime string, pageCount int) error {
	sum := sha256.Sum256([]byte(text))
	params := db.CompleteDocumentExtractionParams{
		DocumentID: documentID,
		Text:       pgtype.Text{String: text, Valid: true},
		Provider:   pgtype.Text{String: provider, Valid: provider != ""},
		Mime:       pgtype.Text{String: mime, Valid: mime != ""},
		CharCount:  pgtype.Int4{Int32: int32(utf8.RuneCountInString(text)), Valid: true},
		TextSha256: pgtype.Text{String: hex.EncodeToString(sum[:]), Valid: true},
	}
	if pageCount > 0 {
		params.PageCount = pgtype.Int4{Int32: int32(pageCount), Valid: true}
	}

	_, err := s.repository.CompleteDocumentExtraction(ctx, params)
	return err
}

func (s *Service) fail(ctx context.Context, documentID uuid.UUID, cause error) error {
	if _, err := s.repository.FailDocumentExtraction(ctx, db.FailDocumentExtractionParams{
		DocumentID:   documentID,
		ErrorMessage: pgtype.Text{String: cause.Error(), Valid: true},
	}); err != nil {
		return err
	}
	return cause
}

func IsPermanent(err error) bool {
	switch {
	case errors.Is(err, ErrUnsupportedFileType),
		errors.Is(err, ErrFileTooLarge),
		errors.Is(err, ErrNoTextExtracted),
		errors.Is(err, ErrEncryptedPDF),
		errors.Is(err, ErrCorruptPDF):
		return true
	}

	var retryabler interface{ Retryable() bool }
	if errors.As(err, &retryabler) {
		return !retryabler.Retryable()
	}

	return false
}

func sniffMIME(data []byte) string {

	ct := http.DetectContentType(data)
	if i := strings.IndexByte(ct, ';'); i >= 0 {
		ct = ct[:i]
	}
	return strings.TrimSpace(ct)
}

var (
	reNewlines   = regexp.MustCompile(`\n{3,}`)
	reSpacesTabs = regexp.MustCompile(`[ \t]+`)
)

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = reSpacesTabs.ReplaceAllString(s, " ")
	s = reNewlines.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

func isEncryptedPDFErr(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "password")
}
