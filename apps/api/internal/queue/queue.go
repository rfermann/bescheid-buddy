package queue

import (
	"context"
	"fmt"

	"github.com/bored-bots/bored-template/internal/queue/workers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type Queue struct {
	Client *river.Client[pgx.Tx]
}

func New(dbPool *pgxpool.Pool, extractor workers.Extractor, analyzer workers.Analyzer, replyGenerator workers.ReplyGenerator, paperLetterDeliverer workers.PaperLetterDeliverer) (*Queue, error) {
	wrks := river.NewWorkers()
	if err := river.AddWorkerSafely(wrks, workers.NewDocumentDataExtractorWorker(extractor)); err != nil {
		return nil, fmt.Errorf("failed to register document data extractor worker: %w", err)
	}
	if err := river.AddWorkerSafely(wrks, workers.NewDocumentAnalyzerWorker(analyzer)); err != nil {
		return nil, fmt.Errorf("failed to register document analyzer worker: %w", err)
	}
	if err := river.AddWorkerSafely(wrks, workers.NewDocumentReplyGenerationWorker(replyGenerator)); err != nil {
		return nil, fmt.Errorf("failed to register document reply generation worker: %w", err)
	}
	if err := river.AddWorkerSafely(wrks, workers.NewPaperLetterDeliveryWorker(paperLetterDeliverer)); err != nil {
		return nil, fmt.Errorf("failed to register paper letter delivery worker: %w", err)
	}

	riverClient, err := river.NewClient[pgx.Tx](riverpgxv5.New(dbPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 5},
		},
		Workers: wrks,
	})
	if err != nil {
		return nil, err
	}

	return &Queue{
		Client: riverClient,
	}, nil
}

func (q *Queue) Start(ctx context.Context) error {
	return q.Client.Start(ctx)
}

func (q *Queue) Stop(ctx context.Context) error {
	return q.Client.Stop(ctx)
}
