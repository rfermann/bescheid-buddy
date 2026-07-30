package workers

import (
	"context"

	"github.com/bored-bots/bored-template/internal/services/extraction_service"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

type DocumentDataExtractionJobArgs struct {
	UserID uuid.UUID `json:"userId"`
	CaseID uuid.UUID `json:"caseId"`
	FileID uuid.UUID `json:"fileId"`
}

type Extractor interface {
	ExtractDocument(ctx context.Context, userID, caseID, documentID uuid.UUID) error
}

type DocumentDataExtractionWorker struct {
	river.WorkerDefaults[DocumentDataExtractionJobArgs]
	extractor Extractor
}

func (DocumentDataExtractionJobArgs) Kind() string {
	return "document-data-extraction"
}

func (DocumentDataExtractionJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	}
}

func NewDocumentDataExtractorWorker(extractor Extractor) *DocumentDataExtractionWorker {
	return &DocumentDataExtractionWorker{extractor: extractor}
}

func (w *DocumentDataExtractionWorker) Work(ctx context.Context, job *river.Job[DocumentDataExtractionJobArgs]) error {
	err := w.extractor.ExtractDocument(ctx, job.Args.UserID, job.Args.CaseID, job.Args.FileID)
	if err != nil && extraction_service.IsPermanent(err) {
		return river.JobCancel(err)
	}

	return err
}
