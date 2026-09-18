package workers

import (
	"context"

	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

type PaperLetterDeliveryJobArgs struct {
	DeliveryID uuid.UUID `json:"deliveryId"`
}

func (PaperLetterDeliveryJobArgs) Kind() string { return "paper-letter-delivery" }

func (PaperLetterDeliveryJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{UniqueOpts: river.UniqueOpts{ByArgs: true}}
}

type PaperLetterDeliverer interface {
	ProcessDelivery(ctx context.Context, deliveryID uuid.UUID) error
}

type PaperLetterDeliveryWorker struct {
	river.WorkerDefaults[PaperLetterDeliveryJobArgs]
	deliverer PaperLetterDeliverer
}

func NewPaperLetterDeliveryWorker(deliverer PaperLetterDeliverer) *PaperLetterDeliveryWorker {
	return &PaperLetterDeliveryWorker{deliverer: deliverer}
}

func (w *PaperLetterDeliveryWorker) Work(ctx context.Context, job *river.Job[PaperLetterDeliveryJobArgs]) error {
	err := w.deliverer.ProcessDelivery(ctx, job.Args.DeliveryID)
	if err != nil && paper_letters_service.IsPermanentDelivery(err) {
		return river.JobCancel(err)
	}
	return err
}
