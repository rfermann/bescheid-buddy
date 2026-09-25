package paper_letters_handler

import (
	"context"
	"time"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type DeliveryListItemBody struct {
	DeliveryID uuid.UUID  `json:"deliveryId"`
	CaseID     *uuid.UUID `json:"caseId" nullable:"true"`
	Recipient  *string    `json:"recipient" nullable:"true"`
	Topic      *string    `json:"topic" nullable:"true"`
	Status     string     `json:"status"`
	SentAt     *time.Time `json:"sentAt" nullable:"true"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type GetPaperLetterDeliveriesOutput struct {
	Body struct {
		Deliveries []DeliveryListItemBody `json:"deliveries"`
	}
}

func (h *Handler) getDeliveries(ctx context.Context, _ *struct{}) (*GetPaperLetterDeliveriesOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	deliveries, err := h.service.GetDeliveries(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get paper letter deliveries", err)
	}

	output := &GetPaperLetterDeliveriesOutput{}
	output.Body.Deliveries = make([]DeliveryListItemBody, len(deliveries))
	for i, delivery := range deliveries {
		item := DeliveryListItemBody{
			DeliveryID: delivery.ID,
			Status:     string(delivery.Status),
			SentAt:     delivery.SentAt,
			UpdatedAt:  delivery.UpdatedAt,
		}
		if delivery.CaseID.Valid {
			caseID := uuid.UUID(delivery.CaseID.Bytes)
			item.CaseID = &caseID
		}
		if delivery.Recipient.Valid {
			item.Recipient = &delivery.Recipient.String
		}
		if delivery.Topic.Valid {
			item.Topic = &delivery.Topic.String
		}
		output.Body.Deliveries[i] = item
	}

	return output, nil
}
