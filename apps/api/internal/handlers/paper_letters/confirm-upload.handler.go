package paper_letters_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ConfirmPaperLetterUploadInput struct {
	UploadID uuid.UUID `path:"uploadId"`
}

type ConfirmPaperLetterUploadOutput struct {
	Body struct {
		UploadID      uuid.UUID `json:"uploadId"`
		Status        string    `json:"status"`
		SizeBytes     int64     `json:"sizeBytes"`
		PageCount     int32     `json:"pageCount"`
		DeliveryPrice int       `json:"deliveryPrice"`
	}
}

func (h *Handler) confirmUpload(ctx context.Context, input *ConfirmPaperLetterUploadInput) (*ConfirmPaperLetterUploadOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	result, err := h.service.ConfirmUpload(ctx, userID, input.UploadID)
	if err != nil {
		switch {
		case errors.Is(err, paper_letters_service.ErrPaperLetterUploadNotFound):
			return nil, huma.Error404NotFound("Paper letter upload not found")
		case errors.Is(err, paper_letters_service.ErrUploadedPDFNotFound):
			return nil, huma.Error400BadRequest("Uploaded PDF was not found in storage")
		case errors.Is(err, paper_letters_service.ErrPDFTooLarge):
			return nil, huma.Error413RequestEntityTooLarge("Uploaded PDF exceeds the maximum allowed size")
		case errors.Is(err, paper_letters_service.ErrInvalidPDF):
			return nil, huma.Error422UnprocessableEntity("Uploaded file is not a valid PDF")
		default:
			return nil, huma.Error500InternalServerError("Failed to confirm paper letter upload", err)
		}
	}
	upload := result.Upload
	output := &ConfirmPaperLetterUploadOutput{}
	output.Body.UploadID = upload.ID
	output.Body.Status = string(upload.Status)
	output.Body.SizeBytes = upload.SizeBytes
	output.Body.PageCount = upload.PageCount
	output.Body.DeliveryPrice = result.DeliveryPrice
	return output, nil
}
