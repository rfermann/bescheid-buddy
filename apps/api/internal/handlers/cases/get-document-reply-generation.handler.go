package cases_handler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetDocumentReplyGenerationInput struct {
	CaseID            uuid.UUID `path:"caseId"`
	DocumentID        uuid.UUID `path:"documentId"`
	ReplyGenerationID uuid.UUID `path:"replyGenerationId"`
}

type DocumentReplyGenerationResultBody struct {
	GermanLetter       string `json:"germanLetter"`
	UserLanguageLetter string `json:"userLanguageLetter"`
	Explanation        string `json:"explanation"`
	Subject            string `json:"subject"`
	Tone               string `json:"tone"`
}

type GetDocumentReplyGenerationOutputBody struct {
	ReplyGenerationID uuid.UUID                          `json:"replyGenerationId"`
	DocumentID        uuid.UUID                          `json:"documentId"`
	Status            string                             `json:"status"`
	JobID             int64                              `json:"jobId,omitempty"`
	Model             string                             `json:"model"`
	Provider          string                             `json:"provider"`
	Language          string                             `json:"language"`
	Instruction       string                             `json:"instruction,omitempty"`
	Result            *DocumentReplyGenerationResultBody `json:"result,omitempty"`
	ErrorMessage      string                             `json:"errorMessage,omitempty"`
	StartedAt         *time.Time                         `json:"startedAt,omitempty"`
	CompletedAt       *time.Time                         `json:"completedAt,omitempty"`
	CreatedAt         time.Time                          `json:"createdAt"`
	UpdatedAt         time.Time                          `json:"updatedAt"`
}

type GetDocumentReplyGenerationOutput struct {
	Body GetDocumentReplyGenerationOutputBody `json:"body"`
}

func (h *Handler) getDocumentReplyGeneration(ctx context.Context, input *GetDocumentReplyGenerationInput) (*GetDocumentReplyGenerationOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	result, err := h.paperLetters.GetSourceGeneration(ctx, userID, input.CaseID, input.DocumentID, input.ReplyGenerationID)
	if err != nil {
		if errors.Is(err, paper_letters_service.ErrGenerationNotFound) {
			return nil, huma.Error404NotFound("Reply generation not found")
		}
		return nil, huma.Error500InternalServerError("Failed to load reply generation", err)
	}
	generation := result.Generation

	body := GetDocumentReplyGenerationOutputBody{
		ReplyGenerationID: generation.ID,
		DocumentID:        input.DocumentID,
		Status:            string(generation.Status),
		JobID:             generation.JobID.Int64,
		Model:             generation.Model,
		Provider:          generation.Provider,
		Language:          generation.Language,
		Instruction:       result.Input.Instruction,
		ErrorMessage:      generation.ErrorMessage.String,
		StartedAt:         generation.StartedAt,
		CompletedAt:       generation.CompletedAt,
		CreatedAt:         generation.CreatedAt,
		UpdatedAt:         generation.UpdatedAt,
	}
	if len(generation.Result) > 0 {
		var result DocumentReplyGenerationResultBody
		if err := json.Unmarshal(generation.Result, &result); err != nil {
			return nil, huma.Error500InternalServerError("Failed to decode reply generation result", err)
		}
		body.Result = &result
	}

	return &GetDocumentReplyGenerationOutput{Body: body}, nil
}
