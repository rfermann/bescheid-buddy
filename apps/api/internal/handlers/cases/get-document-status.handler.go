package cases_handler

import (
	"context"
	"encoding/json"
	"errors"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetDocumentStatusInput struct {
	CaseID     uuid.UUID `path:"caseId"`
	DocumentID uuid.UUID `path:"documentId"`
}

type AnalysisStatusBody struct {
	Status       string    `json:"status"`
	Model        string    `json:"model,omitempty"`
	Provider     string    `json:"provider,omitempty"`
	Text         string    `json:"text,omitempty"`
	Metadata     *Metadata `json:"metadata,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

type Address struct {
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
	Street      string `json:"street,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Country     string `json:"country,omitempty"`
	Category    string `json:"category,omitempty"`
}

type Dates struct {
	LetterDate   string `json:"letter_date,omitempty"`
	DeadlineDate string `json:"deadline_date,omitempty"`
}

type Subject struct {
	Localized string `json:"localized,omitempty"`
	Original  string `json:"original,omitempty"`
}

type Details struct {
	Summary              string   `json:"summary,omitempty"`
	DocumentsRequested   []string `json:"documents_requested,omitempty"`
	Deadlines            []string `json:"deadlines,omitempty"`
	RecommendedNextSteps []string `json:"recommended_next_steps,omitempty"`
	PossibleConsequences []string `json:"possible_consequences,omitempty"`
}

type Metadata struct {
	Receiver        Address `json:"receiver,omitempty"`
	Sender          Address `json:"sender,omitempty"`
	LetterType      string  `json:"letter_type,omitempty"`
	Topic           string  `json:"topic,omitempty"`
	Subject         Subject `json:"subject,omitempty"`
	Dates           Dates   `json:"dates,omitempty"`
	Importance      string  `json:"importance,omitempty"`
	FullTranslation string  `json:"part_1_full_translation,omitempty"`
	Analysis        Details `json:"part_2_ai_analysis,omitempty"`
}

type GetDocumentStatusOutputBody struct {
	DocumentID       uuid.UUID          `json:"documentId"`
	FileName         string             `json:"fileName"`
	UploadStatus     string             `json:"uploadStatus"`
	ExtractionStatus string             `json:"extractionStatus"`
	Analysis         AnalysisStatusBody `json:"analysis"`
}

type GetDocumentStatusOutput struct {
	Body GetDocumentStatusOutputBody `json:"body"`
}

func (h *Handler) getDocumentStatus(ctx context.Context, input *GetDocumentStatusInput) (*GetDocumentStatusOutput, error) {
	if _, ok := appContext.UserIDFromContext(ctx); !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	status, err := h.service.GetDocumentStatus(ctx, input.DocumentID)
	if err != nil {
		if errors.Is(err, cases_service.ErrDocumentNotFound) {
			return nil, huma.Error404NotFound("Document not found")
		}
		return nil, huma.Error500InternalServerError("Failed to load document status", err)
	}

	body := GetDocumentStatusOutputBody{
		DocumentID:       status.DocumentID,
		FileName:         status.FileName,
		UploadStatus:     string(status.UploadStatus),
		ExtractionStatus: "none",
		Analysis:         AnalysisStatusBody{Status: "none"},
	}

	if status.Extraction != nil {
		body.ExtractionStatus = string(status.Extraction.Status)
	}

	if a := status.Analysis; a != nil {
		body.Analysis = AnalysisStatusBody{
			Status:       string(a.Status),
			Model:        a.Model,
			Provider:     a.Provider,
			Text:         a.Text,
			ErrorMessage: a.ErrorMessage,
		}
		if len(a.Metadata) > 0 {
			var meta Metadata
			if err := json.Unmarshal(a.Metadata, &meta); err == nil {
				body.Analysis.Metadata = &meta
			}
		}
	}

	return &GetDocumentStatusOutput{Body: body}, nil
}
