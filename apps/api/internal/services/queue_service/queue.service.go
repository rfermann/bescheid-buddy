package queue_service

import (
	"context"

	"github.com/bored-bots/bored-template/internal/queue/workers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Service struct {
	client *river.Client[pgx.Tx]
}

func New(client *river.Client[pgx.Tx]) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) QueueDocumentDataExtractionJob(ctx context.Context, userID, caseID, fileID uuid.UUID) (int64, error) {
	result, err := s.client.Insert(ctx, workers.DocumentDataExtractionJobArgs{
		UserID: userID,
		CaseID: caseID,
		FileID: fileID,
	}, nil)
	if err != nil {
		return 0, err
	}
	return result.Job.ID, nil
}

func (s *Service) QueueDocumentDataExtractionJobTx(ctx context.Context, tx pgx.Tx, userID, caseID, fileID uuid.UUID) (int64, error) {
	result, err := s.client.InsertTx(ctx, tx, workers.DocumentDataExtractionJobArgs{
		UserID: userID,
		CaseID: caseID,
		FileID: fileID,
	}, nil)
	if err != nil {
		return 0, err
	}
	return result.Job.ID, nil
}

func (s *Service) QueueDocumentAnalysisJobTx(ctx context.Context, tx pgx.Tx, documentID uuid.UUID, model, language string) (int64, error) {
	result, err := s.client.InsertTx(ctx, tx, workers.DocumentAnalysisJobArgs{
		DocumentID: documentID,
		Model:      model,
		Language:   language,
	}, nil)
	if err != nil {
		return 0, err
	}
	return result.Job.ID, nil
}

func (s *Service) QueueDocumentReplyGenerationJobTx(ctx context.Context, tx pgx.Tx, replyGenerationID uuid.UUID) (int64, error) {
	result, err := s.client.InsertTx(ctx, tx, workers.DocumentReplyGenerationJobArgs{
		ReplyGenerationID: replyGenerationID,
	}, nil)
	if err != nil {
		return 0, err
	}
	return result.Job.ID, nil
}
