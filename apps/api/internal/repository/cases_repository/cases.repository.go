package cases_repository

import (
	"context"

	"github.com/bored-bots/bored-template/internal/store"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreateCase(ctx context.Context, params db.CreateCaseParams) (*db.Case, error)
	GetCasesByUserID(ctx context.Context, params db.GetCasesByUserIDParams) ([]*db.Case, error)
	UpdateCaseTitle(ctx context.Context, params db.UpdateCaseTitleParams) error
	CreateCaseDocument(ctx context.Context, params db.CreateCaseDocumentParams) (*db.CaseDocument, error)
	GetCaseDocumentByID(ctx context.Context, params db.GetCaseDocumentByIDParams) (*db.CaseDocument, error)
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
