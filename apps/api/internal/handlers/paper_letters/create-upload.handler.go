package paper_letters_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreatePaperLetterUploadInput struct {
	Body struct {
		Filename string `json:"filename" required:"true" maxLength:"255"`
	}
}

type CreatePaperLetterUploadOutput struct {
	Body struct {
		UploadID  uuid.UUID `json:"uploadId"`
		UploadURL string    `json:"uploadUrl"`
	}
}

func (h *Handler) createUpload(ctx context.Context, input *CreatePaperLetterUploadInput) (*CreatePaperLetterUploadOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	result, err := h.service.CreateUpload(ctx, userID, input.Body.Filename)
	if err != nil {
		if errors.Is(err, paper_letters_service.ErrInvalidUploadFilename) {
			return nil, huma.Error422UnprocessableEntity("Invalid PDF filename")
		}
		return nil, huma.Error500InternalServerError("Failed to create paper letter upload", err)
	}
	output := &CreatePaperLetterUploadOutput{}
	output.Body.UploadID = result.UploadID
	output.Body.UploadURL = result.UploadURL
	return output, nil
}
