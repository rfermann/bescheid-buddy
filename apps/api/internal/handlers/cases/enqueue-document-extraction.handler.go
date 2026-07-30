package cases_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type EnqueueDocumentExtractionInput struct {
	CaseID     uuid.UUID `path:"caseId"`
	DocumentID uuid.UUID `path:"documentId"`
}

type EnqueueDocumentExtractionOutputBody struct {
	DocumentID uuid.UUID `json:"documentId"`
	JobID      int64     `json:"jobId"`
}

type EnqueueDocumentExtractionOutput struct {
	Body EnqueueDocumentExtractionOutputBody `json:"body"`
}

func (h *Handler) enqueueDocumentExtraction(ctx context.Context, input *EnqueueDocumentExtractionInput) (*EnqueueDocumentExtractionOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	result, err := h.service.EnqueueDocumentExtraction(ctx, cases_service.EnqueueExtractionInput{
		UserID:     userID,
		CaseID:     input.CaseID,
		DocumentID: input.DocumentID,
	})
	if err != nil {
		if errors.Is(err, cases_service.ErrDocumentNotUploaded) {
			return nil, huma.Error409Conflict("Document upload has not been confirmed")
		}
		return nil, huma.Error500InternalServerError("Failed to enqueue document extraction job", err)
	}

	return &EnqueueDocumentExtractionOutput{
		Body: EnqueueDocumentExtractionOutputBody{
			DocumentID: result.DocumentID,
			JobID:      result.JobID,
		},
	}, nil
}
