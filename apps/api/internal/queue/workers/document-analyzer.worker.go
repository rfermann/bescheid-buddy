package workers

import (
	"context"

	"github.com/bored-bots/bored-template/internal/services/analysis_service"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

type DocumentAnalysisJobArgs struct {
	DocumentID uuid.UUID `json:"documentId"`
	Model      string    `json:"model"`
	Language   string    `json:"language"`
}

type Analyzer interface {
	AnalyzeDocument(ctx context.Context, documentID uuid.UUID, model, language string) error
}

type DocumentAnalysisWorker struct {
	river.WorkerDefaults[DocumentAnalysisJobArgs]
	analyzer Analyzer
}

func (DocumentAnalysisJobArgs) Kind() string {
	return "document-analysis"
}

func (DocumentAnalysisJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	}
}

func NewDocumentAnalyzerWorker(analyzer Analyzer) *DocumentAnalysisWorker {
	return &DocumentAnalysisWorker{analyzer: analyzer}
}

func (w *DocumentAnalysisWorker) Work(ctx context.Context, job *river.Job[DocumentAnalysisJobArgs]) error {
	err := w.analyzer.AnalyzeDocument(ctx, job.Args.DocumentID, job.Args.Model, job.Args.Language)
	if err != nil && analysis_service.IsPermanent(err) {
		return river.JobCancel(err)
	}

	return err
}
