package cases_repository

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreateCase(ctx context.Context, params db.CreateCaseParams) (*db.Case, error)
	GetCaseByIDAndUserID(ctx context.Context, params db.GetCaseByIDAndUserIDParams) (*db.Case, error)
	GetCasesByUserID(ctx context.Context, params db.GetCasesByUserIDParams) ([]*db.Case, error)
	UpdateCaseTitle(ctx context.Context, params db.UpdateCaseTitleParams) error
	CreateCaseDocument(ctx context.Context, params db.CreateCaseDocumentParams) (*db.CaseDocument, error)
	GetCaseDocumentByID(ctx context.Context, params db.GetCaseDocumentByIDParams) (*db.CaseDocument, error)
	GetCaseDocumentsForDetail(ctx context.Context, params db.GetCaseDocumentsForDetailParams) ([]*db.GetCaseDocumentsForDetailRow, error)
	GetInitialCaseDocument(ctx context.Context, params db.GetInitialCaseDocumentParams) (*db.CaseDocument, error)
	MarkCaseDocumentUploaded(ctx context.Context, params db.MarkCaseDocumentUploadedParams) (*db.CaseDocument, error)
	MarkCaseDocumentUploadFailed(ctx context.Context, params db.MarkCaseDocumentUploadFailedParams) (*db.CaseDocument, error)
	CreateDocumentExtraction(ctx context.Context, params db.CreateDocumentExtractionParams) (*db.CreateDocumentExtractionRow, error)
	GetDocumentExtractionByDocumentID(ctx context.Context, params db.GetDocumentExtractionByDocumentIDParams) (*db.GetDocumentExtractionByDocumentIDRow, error)
	MarkDocumentExtractionProcessing(ctx context.Context, params db.MarkDocumentExtractionProcessingParams) (*db.MarkDocumentExtractionProcessingRow, error)
	CompleteDocumentExtraction(ctx context.Context, params db.CompleteDocumentExtractionParams) (*db.CompleteDocumentExtractionRow, error)
	FailDocumentExtraction(ctx context.Context, params db.FailDocumentExtractionParams) (*db.FailDocumentExtractionRow, error)
	GetDocumentExtractionTextByDocumentID(ctx context.Context, params db.GetDocumentExtractionTextByDocumentIDParams) (*db.GetDocumentExtractionTextByDocumentIDRow, error)
	CreateDocumentAnalysis(ctx context.Context, params db.CreateDocumentAnalysisParams) (*db.CreateDocumentAnalysisRow, error)
	GetDocumentAnalysisByDocumentID(ctx context.Context, params db.GetDocumentAnalysisByDocumentIDParams) (*db.DocumentAnalysisResult, error)
	MarkDocumentAnalysisProcessing(ctx context.Context, params db.MarkDocumentAnalysisProcessingParams) (*db.MarkDocumentAnalysisProcessingRow, error)
	CompleteDocumentAnalysis(ctx context.Context, params db.CompleteDocumentAnalysisParams) (*db.CompleteDocumentAnalysisRow, error)
	FailDocumentAnalysis(ctx context.Context, params db.FailDocumentAnalysisParams) (*db.FailDocumentAnalysisRow, error)
	CreateDocumentReplyGeneration(ctx context.Context, params db.CreateDocumentReplyGenerationParams) (*db.DocumentReplyGeneration, error)
	GetDocumentReplyGenerationByID(ctx context.Context, params db.GetDocumentReplyGenerationByIDParams) (*db.DocumentReplyGeneration, error)
	GetDocumentReplyGenerationByIDAndUserID(ctx context.Context, params db.GetDocumentReplyGenerationByIDAndUserIDParams) (*db.DocumentReplyGeneration, error)
	GetDocumentReplyGenerationByUserAndClientRequestID(ctx context.Context, params db.GetDocumentReplyGenerationByUserAndClientRequestIDParams) (*db.DocumentReplyGeneration, error)
	SetDocumentReplyGenerationJobID(ctx context.Context, params db.SetDocumentReplyGenerationJobIDParams) error
	MarkDocumentReplyGenerationProcessing(ctx context.Context, params db.MarkDocumentReplyGenerationProcessingParams) (*db.DocumentReplyGeneration, error)
	CompleteDocumentReplyGeneration(ctx context.Context, params db.CompleteDocumentReplyGenerationParams) (*db.DocumentReplyGeneration, error)
	FailActiveDocumentReplyGenerationTx(ctx context.Context, tx pgx.Tx, params db.FailActiveDocumentReplyGenerationParams) (*db.FailActiveDocumentReplyGenerationRow, error)
	CreatePaperLetterUpload(ctx context.Context, params db.CreatePaperLetterUploadParams) (*db.PaperLetterUpload, error)
	GetPaperLetterUploadByIDAndUserID(ctx context.Context, params db.GetPaperLetterUploadByIDAndUserIDParams) (*db.PaperLetterUpload, error)
	ConfirmPaperLetterUpload(ctx context.Context, params db.ConfirmPaperLetterUploadParams) (*db.PaperLetterUpload, error)
	CreatePaperLetterDelivery(ctx context.Context, params db.CreatePaperLetterDeliveryParams) (*db.PaperLetterDelivery, error)
	GetPaperLetterDeliveryByID(ctx context.Context, params db.GetPaperLetterDeliveryByIDParams) (*db.PaperLetterDelivery, error)
	GetPaperLetterDeliveryByUserAndClientRequestID(ctx context.Context, params db.GetPaperLetterDeliveryByUserAndClientRequestIDParams) (*db.PaperLetterDelivery, error)
	SetPaperLetterDeliveryJobID(ctx context.Context, params db.SetPaperLetterDeliveryJobIDParams) error
	MarkPaperLetterDeliveryProcessing(ctx context.Context, params db.MarkPaperLetterDeliveryProcessingParams) (*db.PaperLetterDelivery, error)
	MarkPaperLetterDeliverySent(ctx context.Context, params db.MarkPaperLetterDeliverySentParams) (*db.PaperLetterDelivery, error)
	MarkPaperLetterDeliveryManualReview(ctx context.Context, params db.MarkPaperLetterDeliveryManualReviewParams) (*db.PaperLetterDelivery, error)
	FailPaperLetterDelivery(ctx context.Context, params db.FailPaperLetterDeliveryParams) (uuid.UUID, error)
	WithTx(tx pgx.Tx) Repository
}

type CasesRepository struct {
	Store *store.Store
}

func New(store *store.Store) Repository {
	return &CasesRepository{
		Store: store,
	}
}

func (r *CasesRepository) WithTx(tx pgx.Tx) Repository {
	return &CasesRepository{
		Store: &store.Store{
			Pool:    r.Store.Pool,
			Queries: r.Store.Queries.WithTx(tx),
		},
	}
}

func (r *CasesRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.Store.Pool.Begin(ctx)
}

func (r *CasesRepository) CreateCase(ctx context.Context, params db.CreateCaseParams) (*db.Case, error) {
	return r.Store.Queries.CreateCase(ctx, params)
}

func (r *CasesRepository) GetCaseByIDAndUserID(ctx context.Context, params db.GetCaseByIDAndUserIDParams) (*db.Case, error) {
	return r.Store.Queries.GetCaseByIDAndUserID(ctx, params)
}

func (r *CasesRepository) GetCasesByUserID(ctx context.Context, params db.GetCasesByUserIDParams) ([]*db.Case, error) {
	return r.Store.Queries.GetCasesByUserID(ctx, params)
}

func (r *CasesRepository) UpdateCaseTitle(ctx context.Context, params db.UpdateCaseTitleParams) error {
	return r.Store.Queries.UpdateCaseTitle(ctx, params)
}

func (r *CasesRepository) CreateCaseDocument(ctx context.Context, params db.CreateCaseDocumentParams) (*db.CaseDocument, error) {
	return r.Store.Queries.CreateCaseDocument(ctx, params)
}

func (r *CasesRepository) GetCaseDocumentByID(ctx context.Context, params db.GetCaseDocumentByIDParams) (*db.CaseDocument, error) {
	return r.Store.Queries.GetCaseDocumentByID(ctx, params)
}

func (r *CasesRepository) GetCaseDocumentsForDetail(ctx context.Context, params db.GetCaseDocumentsForDetailParams) ([]*db.GetCaseDocumentsForDetailRow, error) {
	return r.Store.Queries.GetCaseDocumentsForDetail(ctx, params)
}

func (r *CasesRepository) GetInitialCaseDocument(ctx context.Context, params db.GetInitialCaseDocumentParams) (*db.CaseDocument, error) {
	return r.Store.Queries.GetInitialCaseDocument(ctx, params)
}

func (r *CasesRepository) MarkCaseDocumentUploaded(ctx context.Context, params db.MarkCaseDocumentUploadedParams) (*db.CaseDocument, error) {
	return r.Store.Queries.MarkCaseDocumentUploaded(ctx, params)
}

func (r *CasesRepository) MarkCaseDocumentUploadFailed(ctx context.Context, params db.MarkCaseDocumentUploadFailedParams) (*db.CaseDocument, error) {
	return r.Store.Queries.MarkCaseDocumentUploadFailed(ctx, params)
}

func (r *CasesRepository) CreateDocumentExtraction(ctx context.Context, params db.CreateDocumentExtractionParams) (*db.CreateDocumentExtractionRow, error) {
	return r.Store.Queries.CreateDocumentExtraction(ctx, params)
}

func (r *CasesRepository) GetDocumentExtractionByDocumentID(ctx context.Context, params db.GetDocumentExtractionByDocumentIDParams) (*db.GetDocumentExtractionByDocumentIDRow, error) {
	return r.Store.Queries.GetDocumentExtractionByDocumentID(ctx, params)
}

func (r *CasesRepository) MarkDocumentExtractionProcessing(ctx context.Context, params db.MarkDocumentExtractionProcessingParams) (*db.MarkDocumentExtractionProcessingRow, error) {
	return r.Store.Queries.MarkDocumentExtractionProcessing(ctx, params)
}

func (r *CasesRepository) CompleteDocumentExtraction(ctx context.Context, params db.CompleteDocumentExtractionParams) (*db.CompleteDocumentExtractionRow, error) {
	return r.Store.Queries.CompleteDocumentExtraction(ctx, params)
}

func (r *CasesRepository) FailDocumentExtraction(ctx context.Context, params db.FailDocumentExtractionParams) (*db.FailDocumentExtractionRow, error) {
	return r.Store.Queries.FailDocumentExtraction(ctx, params)
}

func (r *CasesRepository) GetDocumentExtractionTextByDocumentID(ctx context.Context, params db.GetDocumentExtractionTextByDocumentIDParams) (*db.GetDocumentExtractionTextByDocumentIDRow, error) {
	return r.Store.Queries.GetDocumentExtractionTextByDocumentID(ctx, params)
}

func (r *CasesRepository) CreateDocumentAnalysis(ctx context.Context, params db.CreateDocumentAnalysisParams) (*db.CreateDocumentAnalysisRow, error) {
	return r.Store.Queries.CreateDocumentAnalysis(ctx, params)
}

func (r *CasesRepository) GetDocumentAnalysisByDocumentID(ctx context.Context, params db.GetDocumentAnalysisByDocumentIDParams) (*db.DocumentAnalysisResult, error) {
	return r.Store.Queries.GetDocumentAnalysisByDocumentID(ctx, params)
}

func (r *CasesRepository) MarkDocumentAnalysisProcessing(ctx context.Context, params db.MarkDocumentAnalysisProcessingParams) (*db.MarkDocumentAnalysisProcessingRow, error) {
	return r.Store.Queries.MarkDocumentAnalysisProcessing(ctx, params)
}

func (r *CasesRepository) CompleteDocumentAnalysis(ctx context.Context, params db.CompleteDocumentAnalysisParams) (*db.CompleteDocumentAnalysisRow, error) {
	return r.Store.Queries.CompleteDocumentAnalysis(ctx, params)
}

func (r *CasesRepository) FailDocumentAnalysis(ctx context.Context, params db.FailDocumentAnalysisParams) (*db.FailDocumentAnalysisRow, error) {
	return r.Store.Queries.FailDocumentAnalysis(ctx, params)
}

func (r *CasesRepository) CreateDocumentReplyGeneration(ctx context.Context, params db.CreateDocumentReplyGenerationParams) (*db.DocumentReplyGeneration, error) {
	return r.Store.Queries.CreateDocumentReplyGeneration(ctx, params)
}

func (r *CasesRepository) GetDocumentReplyGenerationByID(ctx context.Context, params db.GetDocumentReplyGenerationByIDParams) (*db.DocumentReplyGeneration, error) {
	return r.Store.Queries.GetDocumentReplyGenerationByID(ctx, params)
}

func (r *CasesRepository) GetDocumentReplyGenerationByIDAndUserID(ctx context.Context, params db.GetDocumentReplyGenerationByIDAndUserIDParams) (*db.DocumentReplyGeneration, error) {
	return r.Store.Queries.GetDocumentReplyGenerationByIDAndUserID(ctx, params)
}

func (r *CasesRepository) GetDocumentReplyGenerationByUserAndClientRequestID(ctx context.Context, params db.GetDocumentReplyGenerationByUserAndClientRequestIDParams) (*db.DocumentReplyGeneration, error) {
	return r.Store.Queries.GetDocumentReplyGenerationByUserAndClientRequestID(ctx, params)
}

func (r *CasesRepository) SetDocumentReplyGenerationJobID(ctx context.Context, params db.SetDocumentReplyGenerationJobIDParams) error {
	return r.Store.Queries.SetDocumentReplyGenerationJobID(ctx, params)
}

func (r *CasesRepository) MarkDocumentReplyGenerationProcessing(ctx context.Context, params db.MarkDocumentReplyGenerationProcessingParams) (*db.DocumentReplyGeneration, error) {
	return r.Store.Queries.MarkDocumentReplyGenerationProcessing(ctx, params)
}

func (r *CasesRepository) CompleteDocumentReplyGeneration(ctx context.Context, params db.CompleteDocumentReplyGenerationParams) (*db.DocumentReplyGeneration, error) {
	return r.Store.Queries.CompleteDocumentReplyGeneration(ctx, params)
}

func (r *CasesRepository) FailActiveDocumentReplyGenerationTx(ctx context.Context, tx pgx.Tx, params db.FailActiveDocumentReplyGenerationParams) (*db.FailActiveDocumentReplyGenerationRow, error) {
	return r.Store.Queries.WithTx(tx).FailActiveDocumentReplyGeneration(ctx, params)
}

func (r *CasesRepository) CreatePaperLetterUpload(ctx context.Context, params db.CreatePaperLetterUploadParams) (*db.PaperLetterUpload, error) {
	return r.Store.Queries.CreatePaperLetterUpload(ctx, params)
}

func (r *CasesRepository) GetPaperLetterUploadByIDAndUserID(ctx context.Context, params db.GetPaperLetterUploadByIDAndUserIDParams) (*db.PaperLetterUpload, error) {
	return r.Store.Queries.GetPaperLetterUploadByIDAndUserID(ctx, params)
}

func (r *CasesRepository) ConfirmPaperLetterUpload(ctx context.Context, params db.ConfirmPaperLetterUploadParams) (*db.PaperLetterUpload, error) {
	return r.Store.Queries.ConfirmPaperLetterUpload(ctx, params)
}

func (r *CasesRepository) CreatePaperLetterDelivery(ctx context.Context, params db.CreatePaperLetterDeliveryParams) (*db.PaperLetterDelivery, error) {
	return r.Store.Queries.CreatePaperLetterDelivery(ctx, params)
}

func (r *CasesRepository) GetPaperLetterDeliveryByID(ctx context.Context, params db.GetPaperLetterDeliveryByIDParams) (*db.PaperLetterDelivery, error) {
	return r.Store.Queries.GetPaperLetterDeliveryByID(ctx, params)
}

func (r *CasesRepository) GetPaperLetterDeliveryByUserAndClientRequestID(ctx context.Context, params db.GetPaperLetterDeliveryByUserAndClientRequestIDParams) (*db.PaperLetterDelivery, error) {
	return r.Store.Queries.GetPaperLetterDeliveryByUserAndClientRequestID(ctx, params)
}

func (r *CasesRepository) SetPaperLetterDeliveryJobID(ctx context.Context, params db.SetPaperLetterDeliveryJobIDParams) error {
	return r.Store.Queries.SetPaperLetterDeliveryJobID(ctx, params)
}

func (r *CasesRepository) MarkPaperLetterDeliveryProcessing(ctx context.Context, params db.MarkPaperLetterDeliveryProcessingParams) (*db.PaperLetterDelivery, error) {
	return r.Store.Queries.MarkPaperLetterDeliveryProcessing(ctx, params)
}

func (r *CasesRepository) MarkPaperLetterDeliverySent(ctx context.Context, params db.MarkPaperLetterDeliverySentParams) (*db.PaperLetterDelivery, error) {
	return r.Store.Queries.MarkPaperLetterDeliverySent(ctx, params)
}

func (r *CasesRepository) MarkPaperLetterDeliveryManualReview(ctx context.Context, params db.MarkPaperLetterDeliveryManualReviewParams) (*db.PaperLetterDelivery, error) {
	return r.Store.Queries.MarkPaperLetterDeliveryManualReview(ctx, params)
}

func (r *CasesRepository) FailPaperLetterDelivery(ctx context.Context, params db.FailPaperLetterDeliveryParams) (uuid.UUID, error) {
	return r.Store.Queries.FailPaperLetterDelivery(ctx, params)
}
