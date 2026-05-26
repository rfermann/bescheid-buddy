package stripe_handler

import (
	"context"
	"fmt"
	"net/http"

	appContext "github.com/bored-bots/bored-template/internal/context"
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

type CreateCheckoutSessionInput struct {
	Body struct {
		StripePriceID string `json:"stripe_price_id" required:"true"`
		PackageName   string `json:"package_name" required:"true"`
	} `json:"body"`
}

type CreateCheckoutSessionOutputBody struct {
	URL string `json:"url"`
}

type CreateCheckoutSessionOutput struct {
	Body CreateCheckoutSessionOutputBody
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "stripe-webhook",
		Method:      http.MethodPost,
		Path:        "/v1/webhooks/stripe",
		Summary:     "Stripe webhook handler",
		Description: "Handles incoming Stripe webhooks",
		Tags:        []string{"Stripe"},
	}, h.handleWebhook)

	huma.Register(api, huma.Operation{
		OperationID: "create-checkout-session",
		Method:      http.MethodPost,
		Path:        "/v1/stripe/checkout-sessions",
		Summary:     "Create Stripe checkout session",
		Description: "Creates a Stripe checkout session for purchasing buddy points",
		Tags:        []string{"Stripe"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.createCheckoutSession)
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

func (h *Handler) createCheckoutSession(ctx context.Context, input *CreateCheckoutSessionInput) (*CreateCheckoutSessionOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	fmt.Println("userID", userID)
	if !ok {
		return nil, huma.Error401Unauthorized(http.StatusText(http.StatusUnauthorized))
	}

	result, err := h.service.CreateCheckoutSession(ctx, stripe_service.CreateCheckoutSessionInput{
		UserID:        userID,
		StripePriceID: input.Body.StripePriceID,
		PackageName:   input.Body.PackageName,
	})
	fmt.Println("result")
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create checkout session", err)
	}

	return &CreateCheckoutSessionOutput{
		Body: CreateCheckoutSessionOutputBody{
			URL: result.URL,
		},
	}, nil
}
