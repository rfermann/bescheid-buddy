package paper_letters_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/bored-bots/bored-template/internal/services/reply_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreateGenerationInput struct {
	Body struct {
		Mode             string              `json:"mode" required:"true" enum:"generate,improve"`
		SourceCaseID     *uuid.UUID          `json:"sourceCaseId,omitempty"`
		SourceDocumentID *uuid.UUID          `json:"sourceDocumentId,omitempty"`
		ClientRequestID  uuid.UUID           `json:"clientRequestId" required:"true"`
		Model            string              `json:"model" required:"true"`
		Locale           string              `json:"locale" required:"true" enum:"en,de,ru,ua,pl,ar"`
		Input            GenerationInputBody `json:"input" required:"true"`
	}
}

type CreateGenerationOutput struct {
	Body GenerationBody
}

func (h *Handler) createGeneration(ctx context.Context, input *CreateGenerationInput) (*CreateGenerationOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	result, err := h.service.Generate(ctx, paper_letters_service.GenerateInput{
		UserID:           userID,
		Mode:             db.DocumentReplyGenerationMode(input.Body.Mode),
		SourceCaseID:     input.Body.SourceCaseID,
		SourceDocumentID: input.Body.SourceDocumentID,
		ClientRequestID:  &input.Body.ClientRequestID,
		Model:            input.Body.Model,
		Locale:           input.Body.Locale,
		Input: reply_service.GenerationInput{
			Sender:              partyInput(input.Body.Input.Sender),
			Recipient:           partyInput(input.Body.Input.Recipient),
			Category:            input.Body.Input.Category,
			Subcategory:         input.Body.Input.Subcategory,
			Subject:             input.Body.Input.Subject,
			Reference:           input.Body.Input.Reference,
			Instruction:         input.Body.Input.Instruction,
			CurrentGermanLetter: input.Body.Input.CurrentGermanLetter,
			Feedback:            input.Body.Input.Feedback,
			Locale:              input.Body.Locale,
		},
	})
	if err != nil {
		return nil, mapGenerationError(err)
	}

	body, err := generationBody(result.GenerationResult)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to encode generation", err)
	}
	return &CreateGenerationOutput{Body: body}, nil
}

func partyInput(value *PartyInputBody) *reply_service.PartyInput {
	if value == nil {
		return nil
	}
	return &reply_service.PartyInput{Name: value.Name, Address: value.Address}
}

func mapGenerationError(err error) error {
	switch {
	case errors.Is(err, paper_letters_service.ErrInvalidSource):
		return huma.Error404NotFound("Paper letter source not found")
	case errors.Is(err, paper_letters_service.ErrSourceNotAnalyzed):
		return huma.Error409Conflict("Paper letter source analysis has not completed")
	case errors.Is(err, paper_letters_service.ErrUnsupportedMode):
		return huma.Error422UnprocessableEntity("Unsupported generation mode")
	case errors.Is(err, paper_letters_service.ErrUnsupportedModel):
		return huma.Error422UnprocessableEntity("Unsupported generation model")
	case errors.Is(err, paper_letters_service.ErrUnsupportedLocale):
		return huma.Error422UnprocessableEntity("Unsupported locale")
	case errors.Is(err, paper_letters_service.ErrInvalidInput):
		return huma.Error422UnprocessableEntity("Invalid generation input")
	case errors.Is(err, paper_letters_service.ErrInsufficientBuddyPoints):
		return huma.Error402PaymentRequired("Insufficient BuddyPoints")
	default:
		return huma.Error500InternalServerError("Failed to create generation", err)
	}
}
