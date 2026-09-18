package cases_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/bored-bots/bored-template/internal/services/reply_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type EnqueueDocumentReplyGenerationInput struct {
	CaseID     uuid.UUID `path:"caseId"`
	DocumentID uuid.UUID `path:"documentId"`
	Body       struct {
		ClientRequestID *uuid.UUID `json:"clientRequestId,omitempty" doc:"Optional client-generated request ID"`
		Model           string     `json:"model" required:"true" doc:"AI model to generate the reply with; see GET /v1/models"`
		Language        string     `json:"language" required:"true" enum:"de,en,ru,ua,pl,ar" doc:"Locale for the reply translation and explanation"`
		Instruction     string     `json:"instruction,omitempty" maxLength:"2500" doc:"Optional instructions for the generated reply"`
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

	result, err := h.paperLetters.Generate(ctx, paper_letters_service.GenerateInput{
		UserID:           userID,
		SourceCaseID:     &input.CaseID,
		SourceDocumentID: &input.DocumentID,
		ClientRequestID:  input.Body.ClientRequestID,
		Model:            input.Body.Model,
		Locale:           input.Body.Language,
		Input: reply_service.GenerationInput{
			Instruction: input.Body.Instruction,
			Locale:      input.Body.Language,
		},
	})
	if err != nil {
		switch {
		case errors.Is(err, paper_letters_service.ErrInvalidSource):
			return nil, huma.Error404NotFound("Document not found")
		case errors.Is(err, paper_letters_service.ErrSourceNotAnalyzed):
			return nil, huma.Error409Conflict("Document analysis has not completed")
		case errors.Is(err, paper_letters_service.ErrUnsupportedModel):
			return nil, huma.Error422UnprocessableEntity("Unsupported reply model")
		case errors.Is(err, paper_letters_service.ErrUnsupportedLocale):
			return nil, huma.Error422UnprocessableEntity("Unsupported language")
		case errors.Is(err, paper_letters_service.ErrInvalidInput):
			return nil, huma.Error422UnprocessableEntity("Invalid generation input")
		case errors.Is(err, paper_letters_service.ErrInsufficientBuddyPoints):
			return nil, huma.Error402PaymentRequired("Insufficient BuddyPoints")
		default:
			return nil, huma.Error500InternalServerError("Failed to enqueue reply generation job", err)
		}
	}

	return &EnqueueDocumentReplyGenerationOutput{
		Body: EnqueueDocumentReplyGenerationOutputBody{
			ReplyGenerationID: result.Generation.ID,
			JobID:             result.Generation.JobID.Int64,
			Status:            string(result.Generation.Status),
		},
	}, nil
}
