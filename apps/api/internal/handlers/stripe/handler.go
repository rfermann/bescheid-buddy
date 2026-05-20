package stripe_handler

import (
	"context"
	"net/http"

	"github.com/bored-bots/bored-template/internal/services/stripe_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stripe/stripe-go/v85"
)

type Handler struct {
	service *stripe_service.Service
}

func NewHandler(service *stripe_service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type WebhookInput struct {
	RawBody         []byte
	StripeSignature string `header:"Stripe-Signature"`
}

type WebhookOutput struct{}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "stripe-webhook",
		Method:      http.MethodPost,
		Path:        "/v1/webhooks/stripe",
		Summary:     "Stripe webhook handler",
		Description: "Handles incoming Stripe webhooks",
		Tags:        []string{"Stripe"},
	}, h.handleWebhook)
}

func (h *Handler) handleWebhook(ctx context.Context, input *WebhookInput) (*WebhookOutput, error) {
	event, err := h.service.VerifyWebhook(input.RawBody, input.StripeSignature)
	if err != nil {
		return nil, huma.Error400BadRequest("Webhook verification failed", err)
	}
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		if err := h.service.HandleCheckoutSessionCompleted(ctx, event); err != nil {
			return nil, huma.Error500InternalServerError("Failed to process checkout session", err)
		}
	default:
		// Unhandled event type
	}

	return &WebhookOutput{}, nil
}
