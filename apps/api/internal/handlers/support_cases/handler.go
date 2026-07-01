package support_cases_handler

import (
	"net/http"

	"github.com/bored-bots/bored-template/internal/services/support_cases_service"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *support_cases_service.Service
}

func NewHandler(service *support_cases_service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-support-case",
		Method:      http.MethodPost,
		Path:        "/support-cases",
		Summary:     "Create a support case",
		Description: "Creates a new support case for the currently logged in user",
		Tags:        []string{"Support Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.createSupportCase)

	huma.Register(api, huma.Operation{
		OperationID: "get-support-cases-by-user",
		Method:      http.MethodGet,
		Path:        "/support-cases",
		Summary:     "Get support cases",
		Description: "Get support cases for the currently logged in user",
		Tags:        []string{"Support Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.GetSupportCasesByUserId)

	huma.Register(api, huma.Operation{
		OperationID: "get-support-case-by-id",
		Method:      http.MethodGet,
		Path:        "/support-cases/{id}",
		Summary:     "Get support case by ID",
		Description: "Get support case by ID",
		Tags:        []string{"Support Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.GetSupportCaseByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-support-case-messages-by-case-id",
		Method:      http.MethodGet,
		Path:        "/support-cases/{caseId}/messages",
		Summary:     "Get support case messages by case ID",
		Description: "Get support case messages by case ID",
		Tags:        []string{"Support Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.GetSupportCaseMessagesBySupportCaseID)
}
