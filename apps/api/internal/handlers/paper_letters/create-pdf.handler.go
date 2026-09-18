package paper_letters_handler

import (
	"context"
	"errors"
	"mime"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_pdf_service"
	"github.com/danielgtaylor/huma/v2"
)

type CreatePDFInput struct {
	Body struct {
		GermanText         string          `json:"germanText" required:"true" maxLength:"24000"`
		Subject            string          `json:"subject" required:"true" maxLength:"500"`
		Sender             *PartyInputBody `json:"sender" required:"true"`
		Recipient          *PartyInputBody `json:"recipient" required:"true"`
		Reference          string          `json:"reference,omitempty" maxLength:"200"`
		Filename           string          `json:"filename,omitempty" maxLength:"200"`
		SignaturePNGBase64 string          `json:"signaturePngBase64,omitempty" maxLength:"700000"`
	}
}

type CreatePDFOutput struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	PageCount          int    `header:"X-Page-Count"`
	Body               []byte
}

func (h *Handler) createPDF(ctx context.Context, input *CreatePDFInput) (*CreatePDFOutput, error) {
	if _, ok := appContext.UserIDFromContext(ctx); !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}
	if input.Body.Sender == nil || input.Body.Recipient == nil {
		return nil, huma.Error422UnprocessableEntity("Invalid paper letter PDF input")
	}

	result, err := h.pdfService.Generate(paper_letters_pdf_service.GenerateInput{
		GermanText: input.Body.GermanText,
		Subject:    input.Body.Subject,
		Sender: paper_letters_pdf_service.Party{
			Name:    input.Body.Sender.Name,
			Address: input.Body.Sender.Address,
		},
		Recipient: paper_letters_pdf_service.Party{
			Name:    input.Body.Recipient.Name,
			Address: input.Body.Recipient.Address,
		},
		Reference:          input.Body.Reference,
		Filename:           input.Body.Filename,
		SignaturePNGBase64: input.Body.SignaturePNGBase64,
	})
	if err != nil {
		switch {
		case errors.Is(err, paper_letters_pdf_service.ErrInvalidInput),
			errors.Is(err, paper_letters_pdf_service.ErrInvalidAddress),
			errors.Is(err, paper_letters_pdf_service.ErrInvalidFilename),
			errors.Is(err, paper_letters_pdf_service.ErrInvalidSignature):
			return nil, huma.Error422UnprocessableEntity("Invalid paper letter PDF input")
		default:
			return nil, huma.Error500InternalServerError("Failed to create paper letter PDF", err)
		}
	}

	return &CreatePDFOutput{
		ContentType:        "application/pdf",
		ContentDisposition: mime.FormatMediaType("attachment", map[string]string{"filename": result.Filename}),
		PageCount:          result.PageCount,
		Body:               result.Bytes,
	}, nil
}
