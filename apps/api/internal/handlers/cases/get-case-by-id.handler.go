package cases_handler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetCaseByIDInput struct {
	CaseID uuid.UUID `path:"caseId"`
}

type GetCaseByIDOutput struct {
	Body struct {
		ID        uuid.UUID          `json:"id"`
		Title     string             `json:"title"`
		CreatedAt time.Time          `json:"createdAt"`
		UpdatedAt time.Time          `json:"updatedAt"`
		Documents []CaseDocumentBody `json:"documents"`
	}
}

type CaseDocumentBody struct {
	DocumentID       uuid.UUID        `json:"documentId"`
	FileName         string           `json:"fileName"`
	UploadStatus     string           `json:"uploadStatus"`
	ExtractionStatus string           `json:"extractionStatus"`
	Analysis         CaseAnalysisBody `json:"analysis"`
}

type CaseAnalysisBody struct {
	Status       string    `json:"status"`
	Text         string    `json:"text,omitempty"`
	Metadata     *Metadata `json:"metadata,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

func (h *Handler) getCaseByID(ctx context.Context, input *GetCaseByIDInput) (*GetCaseByIDOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	caseRecord, err := h.service.GetCaseByID(ctx, userID, input.CaseID)
	if err != nil {
		if errors.Is(err, cases_service.ErrCaseNotFound) {
			return nil, huma.Error404NotFound("Case not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get case", err)
	}

	output := &GetCaseByIDOutput{}
	output.Body.ID = caseRecord.ID
	output.Body.Title = caseRecord.Title.String
	output.Body.CreatedAt = caseRecord.CreatedAt
	output.Body.UpdatedAt = caseRecord.UpdatedAt

	documents, err := h.service.GetCaseDocumentsForDetail(ctx, caseRecord.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get case documents", err)
	}
	output.Body.Documents = make([]CaseDocumentBody, len(documents))
	for i, document := range documents {
		analysisStatus := "none"
		if document.AnalysisStatus.Valid {
			analysisStatus = string(document.AnalysisStatus.DocumentAnalysisStatus)
		}
		extractionStatus := "none"
		if document.ExtractionStatus.Valid {
			extractionStatus = string(document.ExtractionStatus.DocumentAnalysisStatus)
		}
		analysis := CaseAnalysisBody{
			Status:       analysisStatus,
			Text:         document.AnalysisText.String,
			ErrorMessage: document.AnalysisErrorMessage.String,
		}
		if len(document.AnalysisMetadata) > 0 {
			var metadata Metadata
			if json.Unmarshal(document.AnalysisMetadata, &metadata) == nil {
				analysis.Metadata = &metadata
			}
		}
		output.Body.Documents[i] = CaseDocumentBody{
			DocumentID:       document.DocumentID,
			FileName:         document.FileName,
			UploadStatus:     string(document.UploadStatus),
			ExtractionStatus: extractionStatus,
			Analysis:         analysis,
		}
	}

	return output, nil
}
