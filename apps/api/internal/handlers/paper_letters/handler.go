package paper_letters_handler

import (
	"net/http"

	"github.com/bored-bots/bored-template/internal/services/paper_letters_pdf_service"
	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service    *paper_letters_service.Service
	pdfService *paper_letters_pdf_service.Service
}

func NewHandler(service *paper_letters_service.Service, pdfService *paper_letters_pdf_service.Service) *Handler {
	return &Handler{service: service, pdfService: pdfService}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-paper-letter-delivery", Method: http.MethodPost,
		Path: "/paper-letters/deliveries", Summary: "Create a paper letter delivery",
		Description: "Charge BuddyPoints and enqueue an idempotent Pingen delivery",
		Tags:        []string{"Paper Letters"}, Security: []map[string][]string{{"session-cookie": {}}},
		DefaultStatus: http.StatusAccepted,
	}, h.createDelivery)

	huma.Register(api, huma.Operation{
		OperationID: "get-paper-letter-delivery", Method: http.MethodGet,
		Path: "/paper-letters/deliveries/{deliveryId}", Summary: "Get a paper letter delivery",
		Description: "Get the status of an owned paper letter delivery",
		Tags:        []string{"Paper Letters"}, Security: []map[string][]string{{"session-cookie": {}}},
	}, h.getDelivery)

	huma.Register(api, huma.Operation{
		OperationID:   "create-paper-letter-upload",
		Method:        http.MethodPost,
		Path:          "/paper-letters/uploads",
		Summary:       "Create a paper letter upload",
		Description:   "Create an owned pending upload and return a 15-minute presigned PUT URL",
		Tags:          []string{"Paper Letters"},
		Security:      []map[string][]string{{"session-cookie": {}}},
		DefaultStatus: http.StatusCreated,
	}, h.createUpload)

	huma.Register(api, huma.Operation{
		OperationID: "confirm-paper-letter-upload",
		Method:      http.MethodPost,
		Path:        "/paper-letters/uploads/{uploadId}/confirm",
		Summary:     "Confirm a paper letter upload",
		Description: "Validate an owned uploaded PDF and persist its authoritative size, hash, page count, and delivery price",
		Tags:        []string{"Paper Letters"},
		Security:    []map[string][]string{{"session-cookie": {}}},
		Responses: map[string]*huma.Response{
			"413": {Description: "Uploaded PDF exceeds the configured size limit"},
		},
	}, h.confirmUpload)

	huma.Register(api, huma.Operation{
		OperationID: "create-paper-letter-pdf",
		Method:      http.MethodPost,
		Path:        "/paper-letters/pdfs",
		Summary:     "Create a paper letter PDF",
		Description: "Render and validate a downloadable Unicode A4 PDF without charging BuddyPoints",
		Tags:        []string{"Paper Letters"},
		Security:    []map[string][]string{{"session-cookie": {}}},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Generated paper letter PDF",
				Content: map[string]*huma.MediaType{
					"application/pdf": {Schema: &huma.Schema{Type: "string", Format: "binary"}},
				},
			},
		},
	}, h.createPDF)

	huma.Register(api, huma.Operation{
		OperationID:   "create-paper-letter-generation",
		Method:        http.MethodPost,
		Path:          "/paper-letters/generations",
		Summary:       "Create paper letter generation",
		Description:   "Charge BuddyPoints and enqueue an idempotent paper letter generation",
		Tags:          []string{"Paper Letters"},
		Security:      []map[string][]string{{"session-cookie": {}}},
		DefaultStatus: http.StatusAccepted,
	}, h.createGeneration)

	huma.Register(api, huma.Operation{
		OperationID: "get-paper-letter-generation",
		Method:      http.MethodGet,
		Path:        "/paper-letters/generations/{generationId}",
		Summary:     "Get paper letter generation",
		Description: "Get the status and result of an owned paper letter generation",
		Tags:        []string{"Paper Letters"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.getGeneration)
}
