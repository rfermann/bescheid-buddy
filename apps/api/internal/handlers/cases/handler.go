package cases_handler

import (
	"net/http"

	"github.com/bored-bots/bored-template/internal/services/cases_service"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service      *cases_service.Service
	paperLetters *paper_letters_service.Service
}

func NewHandler(service *cases_service.Service, paperLetters *paper_letters_service.Service) *Handler {
	return &Handler{
		service:      service,
		paperLetters: paperLetters,
	}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-case",
		Method:      http.MethodPost,
		Path:        "/cases",
		Summary:     "Create case",
		Description: "Create a new case for document analysis",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.createCase)

	huma.Register(api, huma.Operation{
		OperationID: "get-cases-by-user",
		Method:      http.MethodGet,
		Path:        "/cases",
		Summary:     "Get cases",
		Description: "Get cases for the currently logged in user, newest first",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.getCasesByUserID)

	huma.Register(api, huma.Operation{
		OperationID: "get-case-by-id",
		Method:      http.MethodGet,
		Path:        "/cases/{caseId}",
		Summary:     "Get case",
		Description: "Get a case belonging to the currently logged in user",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.getCaseByID)

	huma.Register(api, huma.Operation{
		OperationID: "generate-upload-url",
		Method:      http.MethodPost,
		Path:        "/cases/{caseId}/documents/upload-url",
		Summary:     "Generate upload URL",
		Description: "Generate a signed URL for uploading a document to a case",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.generateUploadURL)

	huma.Register(api, huma.Operation{
		OperationID: "confirm-document-upload",
		Method:      http.MethodPost,
		Path:        "/cases/{caseId}/documents/{documentId}/confirm",
		Summary:     "Confirm document upload",
		Description: "Verify an uploaded document landed in storage and mark it uploaded",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.confirmUpload)

	huma.Register(api, huma.Operation{
		OperationID: "enqueue-document-extraction",
		Method:      http.MethodPost,
		Path:        "/cases/{caseId}/documents/{documentId}/extract",
		Summary:     "Enqueue document extraction",
		Description: "Enqueue a River job to extract data from an uploaded document",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.enqueueDocumentExtraction)

	huma.Register(api, huma.Operation{
		OperationID: "enqueue-document-analysis",
		Method:      http.MethodPost,
		Path:        "/cases/{caseId}/documents/{documentId}/analyse",
		Summary:     "Enqueue document analysis",
		Description: "Enqueue a River job to analyze an extracted document with the chosen AI model",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.enqueueDocumentAnalysis)

	huma.Register(api, huma.Operation{
		OperationID: "get-document-status",
		Method:      http.MethodGet,
		Path:        "/cases/{caseId}/documents/{documentId}",
		Summary:     "Get document status",
		Description: "Get a document's upload, extraction and analysis status, including the analysis result once complete",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.getDocumentStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "enqueue-document-reply-generation",
		Method:        http.MethodPost,
		Path:          "/cases/{caseId}/documents/{documentId}/reply-generations",
		Summary:       "Enqueue document reply generation",
		Description:   "Enqueue a River job to generate a reply to an analyzed document",
		Tags:          []string{"Cases"},
		Security:      []map[string][]string{{"session-cookie": {}}},
		DefaultStatus: http.StatusAccepted,
	}, h.enqueueDocumentReplyGeneration)

	huma.Register(api, huma.Operation{
		OperationID: "get-document-reply-generation",
		Method:      http.MethodGet,
		Path:        "/cases/{caseId}/documents/{documentId}/reply-generations/{replyGenerationId}",
		Summary:     "Get document reply generation",
		Description: "Get the status and result of a document reply generation",
		Tags:        []string{"Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.getDocumentReplyGeneration)
}
