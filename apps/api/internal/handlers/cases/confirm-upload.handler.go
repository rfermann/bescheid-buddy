package cases_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ConfirmUploadInput struct {
	CaseID     uuid.UUID `path:"caseId"`
	DocumentID uuid.UUID `path:"documentId"`
}

type ConfirmUploadOutputBody struct {
	DocumentID   uuid.UUID `json:"documentId"`
	UploadStatus string    `json:"uploadStatus"`
}

type ConfirmUploadOutput struct {
	Body ConfirmUploadOutputBody `json:"body"`
}

func (h *Handler) confirmUpload(ctx context.Context, input *ConfirmUploadInput) (*ConfirmUploadOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	document, err := h.service.ConfirmDocumentUpload(ctx, cases_service.ConfirmDocumentUploadInput{
		UserID:     userID,
		CaseID:     input.CaseID,
		DocumentID: input.DocumentID,
	})
	if err != nil {
		if errors.Is(err, cases_service.ErrUploadNotFound) {
			return nil, huma.Error400BadRequest("Uploaded file was not found in storage")
		}
		if errors.Is(err, cases_service.ErrFileTooLarge) {
			return nil, huma.Error413RequestEntityTooLarge("Uploaded file exceeds the maximum allowed size")
		}
		return nil, huma.Error500InternalServerError("Failed to confirm upload", err)
	}

	return &ConfirmUploadOutput{
		Body: ConfirmUploadOutputBody{
			DocumentID:   document.ID,
			UploadStatus: string(document.UploadStatus),
		},
	}, nil
}
