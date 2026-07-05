package support_cases_handler

import (
	"context"
	"fmt"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreateSupportCaseMessageInputBody struct {
	Message       string  `json:"message" required:"true" maxLength:"255"`
	IsInternal    *bool   `json:"isInternal" default:"false"`
	IsUserMessage *bool   `json:"isUserMessage" default:"true"`
}

type CreateSupportCaseMessageInput struct {
	CaseID uuid.UUID                         `path:"caseId"`
	Body   CreateSupportCaseMessageInputBody `json:"body"`
}

type CreateSupportCaseMessageOutputBody struct {
	Message *db.SupportCaseMessage `json:"message"`
}

type CreateSupportCaseMessageOutput struct {
	Body CreateSupportCaseMessageOutputBody `json:"body"`
}

func (h *Handler) CreateSupportCaseMessage(ctx context.Context, input *CreateSupportCaseMessageInput) (*CreateSupportCaseMessageOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	fmt.Println("Creating support case message", input.Body.IsUserMessage)

	isUserMessage := true
	if input.Body.IsUserMessage != nil {
		isUserMessage = *input.Body.IsUserMessage
	}

	isInternal := false
	if input.Body.IsInternal != nil {
		isInternal = *input.Body.IsInternal
	}

	params := support_cases_service.CreateSupportCaseMessageInput{
		CaseID:        input.CaseID,
		Message:       input.Body.Message,
		UserID:        userID,
		IsUserMessage: isUserMessage,
		IsInternal:    isInternal,
	}

	message, err := h.service.CreateSupportCaseMessage(ctx, params)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create support case message", err)
	}

	return &CreateSupportCaseMessageOutput{
		Body: CreateSupportCaseMessageOutputBody{
			Message: message,
		},
	}, nil
}
