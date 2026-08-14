package workers

import (
	"context"

	"github.com/bored-bots/bored-template/internal/services/reply_service"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

type DocumentReplyGenerationJobArgs struct {
	ReplyGenerationID uuid.UUID `json:"replyGenerationId"`
}

func (DocumentReplyGenerationJobArgs) Kind() string {
	return "document-reply-generation"
}

func (DocumentReplyGenerationJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	}
}

type ReplyGenerator interface {
	Generate(ctx context.Context, replyGenerationID uuid.UUID) error
}

type DocumentReplyGenerationWorker struct {
	river.WorkerDefaults[DocumentReplyGenerationJobArgs]
	generator ReplyGenerator
}

func NewDocumentReplyGenerationWorker(generator ReplyGenerator) *DocumentReplyGenerationWorker {
	return &DocumentReplyGenerationWorker{generator: generator}
}

func (w *DocumentReplyGenerationWorker) Work(ctx context.Context, job *river.Job[DocumentReplyGenerationJobArgs]) error {
	err := w.generator.Generate(ctx, job.Args.ReplyGenerationID)
	if err != nil && reply_service.IsPermanent(err) {
		return river.JobCancel(err)
	}
	return err
}
