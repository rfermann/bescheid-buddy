package cases_service

import (
	"context"
	"errors"

	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/services/document_service"
	"github.com/bored-bots/bored-template/internal/services/llm_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrUploadNotFound       = errors.New("uploaded file was not found in storage")
	ErrDocumentNotUploaded  = errors.New("document upload has not been confirmed")
	ErrFileTooLarge         = errors.New("uploaded file exceeds the maximum allowed size")
	ErrDocumentNotExtracted = errors.New("document extraction has not completed")
	ErrDocumentNotFound     = errors.New("document was not found")
)

type DocumentStorage interface {
	PresignUpload(ctx context.Context, userID, caseID, documentID uuid.UUID) (string, error)
	StatDocumentSize(ctx context.Context, userID, caseID, documentID uuid.UUID) (int64, error)
	RemoveDocument(ctx context.Context, userID, caseID, documentID uuid.UUID) error
}

type DocumentJobQueue interface {
	QueueDocumentDataExtractionJobTx(ctx context.Context, tx pgx.Tx, userID, caseID, fileID uuid.UUID) (int64, error)
	QueueDocumentAnalysisJobTx(ctx context.Context, tx pgx.Tx, documentID uuid.UUID, model, language string) (int64, error)
}

type Service struct {
	repository cases_repository.Repository
	storage    DocumentStorage
	queue      DocumentJobQueue
	maxBytes   int64
}

func New(repository cases_repository.Repository, storage DocumentStorage, queue DocumentJobQueue, maxBytes int64) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
		queue:      queue,
		maxBytes:   maxBytes,
	}
}

type CreateCaseInput struct {
	UserID uuid.UUID
	Title  string
}

func (s *Service) CreateCase(ctx context.Context, input CreateCaseInput) (*db.Case, error) {
	return s.repository.CreateCase(ctx, db.CreateCaseParams{
		UserID: input.UserID,
		Title:  pgtype.Text{String: input.Title, Valid: input.Title != ""},
	})
}

type CreateDocumentUploadInput struct {
	UserID   uuid.UUID
	CaseID   uuid.UUID
	FileName string
}

type CreateDocumentUploadResult struct {
	UploadURL  string
	DocumentID uuid.UUID
}

func (s *Service) CreateDocumentUpload(ctx context.Context, input CreateDocumentUploadInput) (*CreateDocumentUploadResult, error) {
	document, err := s.repository.CreateCaseDocument(ctx, db.CreateCaseDocumentParams{
		CaseID:   input.CaseID,
		FileName: input.FileName,
	})
	if err != nil {
		return nil, err
	}

	uploadURL, err := s.storage.PresignUpload(ctx, input.UserID, input.CaseID, document.ID)
	if err != nil {
		return nil, err
	}

	return &CreateDocumentUploadResult{
		UploadURL:  uploadURL,
		DocumentID: document.ID,
	}, nil
}

type ConfirmDocumentUploadInput struct {
	UserID     uuid.UUID
	CaseID     uuid.UUID
	DocumentID uuid.UUID
}

func (s *Service) ConfirmDocumentUpload(ctx context.Context, input ConfirmDocumentUploadInput) (*db.CaseDocument, error) {
	size, err := s.storage.StatDocumentSize(ctx, input.UserID, input.CaseID, input.DocumentID)
	if err != nil {
		if errors.Is(err, document_service.ErrObjectNotFound) {
			if _, ferr := s.repository.MarkCaseDocumentUploadFailed(ctx, db.MarkCaseDocumentUploadFailedParams{ID: input.DocumentID}); ferr != nil {
				return nil, ferr
			}
			return nil, ErrUploadNotFound
		}
		return nil, err
	}

	if size > s.maxBytes {
		_ = s.storage.RemoveDocument(ctx, input.UserID, input.CaseID, input.DocumentID)
		if _, ferr := s.repository.MarkCaseDocumentUploadFailed(ctx, db.MarkCaseDocumentUploadFailedParams{ID: input.DocumentID}); ferr != nil {
			return nil, ferr
		}
		return nil, ErrFileTooLarge
	}

	return s.repository.MarkCaseDocumentUploaded(ctx, db.MarkCaseDocumentUploadedParams{ID: input.DocumentID})
}

type EnqueueExtractionInput struct {
	UserID     uuid.UUID
	CaseID     uuid.UUID
	DocumentID uuid.UUID
}

type EnqueueExtractionResult struct {
	DocumentID uuid.UUID
	JobID      int64
}

func (s *Service) EnqueueDocumentExtraction(ctx context.Context, input EnqueueExtractionInput) (*EnqueueExtractionResult, error) {
	document, err := s.repository.GetCaseDocumentByID(ctx, db.GetCaseDocumentByIDParams{ID: input.DocumentID})
	if err != nil {
		return nil, err
	}

	if document.UploadStatus != db.DocumentUploadStatusUploaded {
		return nil, ErrDocumentNotUploaded
	}

	existing, err := s.repository.GetDocumentExtractionByDocumentID(ctx, db.GetDocumentExtractionByDocumentIDParams{DocumentID: input.DocumentID})
	if err == nil {
		return &EnqueueExtractionResult{DocumentID: document.ID, JobID: existing.JobID.Int64}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	txRepo := s.repository.WithTx(tx)

	jobID, err := s.queue.QueueDocumentDataExtractionJobTx(ctx, tx, input.UserID, input.CaseID, input.DocumentID)
	if err != nil {
		return nil, err
	}

	if _, err := txRepo.CreateDocumentExtraction(ctx, db.CreateDocumentExtractionParams{
		DocumentID: input.DocumentID,
		JobID:      pgtype.Int8{Int64: jobID, Valid: true},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &EnqueueExtractionResult{
		DocumentID: document.ID,
		JobID:      jobID,
	}, nil
}

type EnqueueAnalysisInput struct {
	CaseID     uuid.UUID
	DocumentID uuid.UUID
	Model      llm_service.Model
	Language   string
}

type EnqueueAnalysisResult struct {
	DocumentID uuid.UUID
	JobID      int64
}

func (s *Service) EnqueueDocumentAnalysis(ctx context.Context, input EnqueueAnalysisInput) (*EnqueueAnalysisResult, error) {
	extraction, err := s.repository.GetDocumentExtractionByDocumentID(ctx, db.GetDocumentExtractionByDocumentIDParams{DocumentID: input.DocumentID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDocumentNotExtracted
		}
		return nil, err
	}
	if extraction.Status != db.DocumentAnalysisStatusCompleted {
		return nil, ErrDocumentNotExtracted
	}

	existing, err := s.repository.GetDocumentAnalysisByDocumentID(ctx, db.GetDocumentAnalysisByDocumentIDParams{DocumentID: input.DocumentID})
	if err == nil {
		return &EnqueueAnalysisResult{DocumentID: input.DocumentID, JobID: existing.JobID.Int64}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	txRepo := s.repository.WithTx(tx)

	jobID, err := s.queue.QueueDocumentAnalysisJobTx(ctx, tx, input.DocumentID, string(input.Model), input.Language)
	if err != nil {
		return nil, err
	}

	if _, err := txRepo.CreateDocumentAnalysis(ctx, db.CreateDocumentAnalysisParams{
		DocumentID: input.DocumentID,
		JobID:      pgtype.Int8{Int64: jobID, Valid: true},
		Model:      pgtype.Text{String: string(input.Model), Valid: input.Model != ""},
		Language:   pgtype.Text{String: input.Language, Valid: input.Language != ""},
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &EnqueueAnalysisResult{
		DocumentID: input.DocumentID,
		JobID:      jobID,
	}, nil
}

type ExtractionStatus struct {
	Status       db.DocumentAnalysisStatus
	ErrorMessage string
}

type AnalysisStatus struct {
	Status       db.DocumentAnalysisStatus
	Model        string
	Provider     string
	Text         string
	Metadata     []byte
	ErrorMessage string
}

type DocumentStatus struct {
	DocumentID   uuid.UUID
	FileName     string
	UploadStatus db.DocumentUploadStatus
	Extraction   *ExtractionStatus
	Analysis     *AnalysisStatus
}

func (s *Service) GetDocumentStatus(ctx context.Context, documentID uuid.UUID) (*DocumentStatus, error) {
	document, err := s.repository.GetCaseDocumentByID(ctx, db.GetCaseDocumentByIDParams{ID: documentID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	status := &DocumentStatus{
		DocumentID:   document.ID,
		FileName:     document.FileName,
		UploadStatus: document.UploadStatus,
	}

	extraction, err := s.repository.GetDocumentExtractionByDocumentID(ctx, db.GetDocumentExtractionByDocumentIDParams{DocumentID: documentID})
	switch {
	case err == nil:
		status.Extraction = &ExtractionStatus{
			Status:       extraction.Status,
			ErrorMessage: extraction.ErrorMessage.String,
		}
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}

	analysis, err := s.repository.GetDocumentAnalysisByDocumentID(ctx, db.GetDocumentAnalysisByDocumentIDParams{DocumentID: documentID})
	switch {
	case err == nil:
		status.Analysis = &AnalysisStatus{
			Status:       analysis.Status,
			Model:        analysis.Model.String,
			Provider:     analysis.Provider.String,
			Text:         analysis.AnalysisText.String,
			Metadata:     analysis.AnalysisMetadata,
			ErrorMessage: analysis.ErrorMessage.String,
		}
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}

	return status, nil
}
