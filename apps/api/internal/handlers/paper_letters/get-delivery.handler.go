package paper_letters_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetPaperLetterDeliveryInput struct {
	DeliveryID uuid.UUID `path:"deliveryId"`
}

type GetPaperLetterDeliveryOutput struct{ Body DeliveryBody }

func (h *Handler) getDelivery(ctx context.Context, input *GetPaperLetterDeliveryInput) (*GetPaperLetterDeliveryOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	delivery, err := h.service.GetDelivery(ctx, userID, input.DeliveryID)
	if err != nil {
		if errors.Is(err, paper_letters_service.ErrDeliveryNotFound) {
			return nil, huma.Error404NotFound("Paper letter delivery not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get paper letter delivery", err)
	}
	return &GetPaperLetterDeliveryOutput{Body: deliveryBody(delivery)}, nil
}
