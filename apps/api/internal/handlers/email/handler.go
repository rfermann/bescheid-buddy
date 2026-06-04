package email_handler

import (
	"context"
	"net/http"

	"github.com/bored-bots/bored-template/internal/services/email_service"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *email_service.Service
}

func NewHandler(service *email_service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type SendEmailInputBody struct {
	HTML     string `json:"html" example:"<h1>Hello</h1>"`
	Subject  string `json:"subject" example:"Welcome"`
	Sender   string `json:"sender" example:"noreply@example.com"`
	Receiver string `json:"receiver" example:"user@example.com"`
}

type SendEmailInput struct {
	Body SendEmailInputBody `json:"body"`
}

type SendEmailOutputBody struct {
	Message string `json:"message" example:"Email logged to console"`
}

type SendEmailOutput struct {
	Body SendEmailOutputBody `json:"body"`
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "send-email",
		Method:      http.MethodPost,
		Path:        "/v1/email/send",
		Summary:     "Send email",
		Description: "Receives email details and logs them to console",
		Tags:        []string{"Email"},
		Security: []map[string][]string{
			{"apiKeyAuth": {}},
		},
	}, h.sendEmail)
}

func (h *Handler) sendEmail(ctx context.Context, input *SendEmailInput) (*SendEmailOutput, error) {
	err := h.service.SendEmail(ctx, email_service.SendEmailInput{
		HTML:     input.Body.HTML,
		Receiver: input.Body.Receiver,
		Sender:   input.Body.Sender,
		Subject:  input.Body.Subject,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to send email", err)
	}

	return &SendEmailOutput{
		Body: SendEmailOutputBody{
			Message: "Email sent successfully",
		},
	}, nil
}
