package cases_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/bored-bots/bored-template/internal/services/llm_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type EnqueueDocumentReplyGenerationInput struct {
	CaseID     uuid.UUID `path:"caseId"`
	DocumentID uuid.UUID `path:"documentId"`
	Body       struct {
		Model       string `json:"model" required:"true" doc:"AI model to generate the reply with; see GET /v1/models"`
		Language    string `json:"language" required:"true" enum:"de,en,ru,uk,ua,pl,ar" doc:"Locale for the reply translation and explanation"`
		Instruction string `json:"instruction,omitempty" maxLength:"2500" doc:"Optional instructions for the generated reply"`
	}
}

type EnqueueDocumentReplyGenerationOutputBody struct {
	ReplyGenerationID uuid.UUID `json:"replyGenerationId"`
	JobID             int64     `json:"jobId"`
	Status            string    `json:"status"`
}

type EnqueueDocumentReplyGenerationOutput struct {
	Body EnqueueDocumentReplyGenerationOutputBody `json:"body"`
}

func (h *Handler) enqueueDocumentReplyGeneration(ctx context.Context, input *EnqueueDocumentReplyGenerationInput) (*EnqueueDocumentReplyGenerationOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	model, err := llm_service.ParseModel(input.Body.Model)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Unsupported reply model")
	}

	result, err := h.service.EnqueueDocumentReplyGeneration(ctx, cases_service.EnqueueReplyGenerationInput{
		UserID:      userID,
		CaseID:      input.CaseID,
		DocumentID:  input.DocumentID,
		Model:       model,
		Language:    input.Body.Language,
		Instruction: input.Body.Instruction,
	})
	if err != nil {
		switch {
		case errors.Is(err, cases_service.ErrDocumentNotFound):
			return nil, huma.Error404NotFound("Document not found")
		case errors.Is(err, cases_service.ErrDocumentNotAnalyzed):
			return nil, huma.Error409Conflict("Document analysis has not completed")
		case errors.Is(err, cases_service.ErrUnsupportedLanguage):
			return nil, huma.Error422UnprocessableEntity("Unsupported language")
		default:
			return nil, huma.Error500InternalServerError("Failed to enqueue reply generation job", err)
		}
	}

	return &EnqueueDocumentReplyGenerationOutput{
		Body: EnqueueDocumentReplyGenerationOutputBody{
			ReplyGenerationID: result.ReplyGenerationID,
			JobID:             result.JobID,
			Status:            "pending",
		},
	}, nil
}
