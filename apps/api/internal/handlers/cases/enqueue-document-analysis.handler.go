package cases_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type EnqueueDocumentAnalysisInput struct {
	CaseID     uuid.UUID `path:"caseId"`
	DocumentID uuid.UUID `path:"documentId"`
	Body       struct {
		Model    string `json:"model" doc:"AI model to run the analysis with; see GET /v1/models"`
		Language string `json:"language,omitempty" doc:"Locale to write the analysis output in (e.g. de, en, ru, ua, pl, ar)"`
	}
}

type EnqueueDocumentAnalysisOutputBody struct {
	DocumentID uuid.UUID `json:"documentId"`
	JobID      int64     `json:"jobId"`
}

type EnqueueDocumentAnalysisOutput struct {
	Body EnqueueDocumentAnalysisOutputBody `json:"body"`
}

func (h *Handler) enqueueDocumentAnalysis(ctx context.Context, input *EnqueueDocumentAnalysisInput) (*EnqueueDocumentAnalysisOutput, error) {
	if _, ok := appContext.UserIDFromContext(ctx); !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	result, err := h.service.EnqueueDocumentAnalysis(ctx, cases_service.EnqueueAnalysisInput{
		CaseID:     input.CaseID,
		DocumentID: input.DocumentID,
		Model:      input.Body.Model,
		Language:   input.Body.Language,
	})
	if err != nil {
		switch {
		case errors.Is(err, cases_service.ErrUnsupportedModel):
			return nil, huma.Error422UnprocessableEntity("Unsupported analysis model")
		case errors.Is(err, cases_service.ErrDocumentNotExtracted):
			return nil, huma.Error409Conflict("Document extraction has not completed")
		}
		return nil, huma.Error500InternalServerError("Failed to enqueue document analysis job", err)
	}

	return &EnqueueDocumentAnalysisOutput{
		Body: EnqueueDocumentAnalysisOutputBody{
			DocumentID: result.DocumentID,
			JobID:      result.JobID,
		},
	}, nil
}
