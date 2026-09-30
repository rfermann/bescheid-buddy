package paper_letters_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetPaperLetterDeliveryPDFURLsInput struct {
	DeliveryID uuid.UUID `path:"deliveryId"`
}

type GetPaperLetterDeliveryPDFURLsOutput struct {
	Body struct {
		InlineURL   string `json:"inlineUrl" format:"uri"`
		DownloadURL string `json:"downloadUrl" format:"uri"`
	}
}

func (h *Handler) getDeliveryPDFURLs(ctx context.Context, input *GetPaperLetterDeliveryPDFURLsInput) (*GetPaperLetterDeliveryPDFURLsOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	urls, err := h.service.GetDeliveryPDFURLs(ctx, userID, input.DeliveryID)
	if err != nil {
		if errors.Is(err, paper_letters_service.ErrDeliveryNotFound) {
			return nil, huma.Error404NotFound("Paper letter delivery not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get paper letter delivery PDF URLs", err)
	}
	output := &GetPaperLetterDeliveryPDFURLsOutput{}
	output.Body.InlineURL = urls.InlineURL
	output.Body.DownloadURL = urls.DownloadURL
	return output, nil
}
