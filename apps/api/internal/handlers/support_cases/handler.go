package support_cases_handler

import (
	"net/http"

	"github.com/bored-bots/bored-template/internal/middleware"
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

	huma.Register(api, huma.Operation{
		OperationID: "create-support-case-message",
		Method:      http.MethodPost,
		Path:        "/support-cases/{caseId}/messages",
		Summary:     "Create support case message",
		Description: "Create a new message for a support case",
		Tags:        []string{"Support Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.CreateSupportCaseMessage)

	huma.Register(api, huma.Operation{
		OperationID: "get-all-support-cases",
		Method:      http.MethodGet,
		Path:        "/admin/support-cases",
		Summary:     "Get all support cases",
		Description: "Get all support cases (admin only)",
		Tags:        []string{"Support Cases"},
		Security:    []map[string][]string{{"session-cookie": {middleware.RoleAdmin}}},
	}, h.GetAllSupportCases)

	huma.Register(api, huma.Operation{
		OperationID: "update-support-case-status",
		Method:      http.MethodPatch,
		Path:        "/support-cases/{id}/status",
		Summary:     "Update support case status",
		Description: "Update the status of a support case",
		Tags:        []string{"Support Cases"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.UpdateSupportCaseStatus)
}
