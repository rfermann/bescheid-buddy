package cases_handler

import (
	"context"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GenerateUploadURLInputBody struct {
	FileName string `json:"fileName" required:"true" maxLength:"255"`
}

type GenerateUploadURLInput struct {
	CaseID uuid.UUID                  `path:"caseId"`
	Body   GenerateUploadURLInputBody `json:"body"`
}

type GenerateUploadURLOutputBody struct {
	UploadURL  string    `json:"uploadUrl"`
	DocumentID uuid.UUID `json:"documentId"`
}

type GenerateUploadURLOutput struct {
	Body GenerateUploadURLOutputBody `json:"body"`
}

func (h *Handler) generateUploadURL(ctx context.Context, input *GenerateUploadURLInput) (*GenerateUploadURLOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	result, err := h.service.CreateDocumentUpload(ctx, cases_service.CreateDocumentUploadInput{
		UserID:   userID,
		CaseID:   input.CaseID,
		FileName: input.Body.FileName,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create upload", err)
	}

	return &GenerateUploadURLOutput{
		Body: GenerateUploadURLOutputBody{
			UploadURL:  result.UploadURL,
			DocumentID: result.DocumentID,
		},
	}, nil
}
