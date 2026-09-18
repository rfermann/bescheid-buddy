package paper_letters_handler

import (
	"context"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreatePaperLetterDeliveryInput struct {
	Body struct {
		UploadID        uuid.UUID `json:"uploadId" required:"true"`
		ClientRequestID uuid.UUID `json:"clientRequestId" required:"true"`
	}
}

type CreatePaperLetterDeliveryOutput struct{ Body DeliveryBody }

func (h *Handler) createDelivery(ctx context.Context, input *CreatePaperLetterDeliveryInput) (*CreatePaperLetterDeliveryOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	result, err := h.service.CreateDelivery(ctx, paper_letters_service.CreateDeliveryInput{
		UserID: userID, UploadID: input.Body.UploadID, ClientRequestID: input.Body.ClientRequestID,
	})
	if err != nil {
		switch {
		case errors.Is(err, paper_letters_service.ErrPaperLetterUploadNotFound):
			return nil, huma.Error404NotFound("Paper letter upload not found")
		case errors.Is(err, paper_letters_service.ErrUploadNotConfirmed):
			return nil, huma.Error409Conflict("Paper letter upload is not confirmed")
		case errors.Is(err, paper_letters_service.ErrInsufficientBuddyPoints):
			return nil, huma.Error402PaymentRequired("Insufficient BuddyPoints")
		case errors.Is(err, paper_letters_service.ErrInvalidInput):
			return nil, huma.Error422UnprocessableEntity("Invalid delivery input")
		default:
			return nil, huma.Error500InternalServerError("Failed to create paper letter delivery", err)
		}
	}
	return &CreatePaperLetterDeliveryOutput{Body: deliveryBody(result.Delivery)}, nil
}
