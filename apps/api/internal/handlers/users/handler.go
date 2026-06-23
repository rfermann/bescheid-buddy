package users_handler

import (
	"context"
	"net/http"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/users_service"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *users_service.Service
}

func NewHandler(service *users_service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type GetBuddyPointsOutputBody struct {
	MonthlyPoints  int32 `json:"monthly_points"`
	LifetimePoints int32 `json:"lifetime_points"`
}

type GetBuddyPointsOutput struct {
	Body GetBuddyPointsOutputBody `json:"body"`
}

type GetBuddyPointsLedgerOutputBody struct {
	Entries []users_service.LedgerEntry `json:"entries"`
}

type GetBuddyPointsLedgerOutput struct {
	Body GetBuddyPointsLedgerOutputBody `json:"body"`
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-buddy-points",
		Method:      http.MethodGet,
		Path:        "/v1/user/buddy-points",
		Summary:     "Get buddy points",
		Description: "Returns the monthly and lifetime buddy points for the currently logged in user",
		Tags:        []string{"User"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.getBuddyPoints)

	huma.Register(api, huma.Operation{
		OperationID: "get-buddy-points-ledger",
		Method:      http.MethodGet,
		Path:        "/v1/user/buddy-points-ledger",
		Summary:     "Get buddy points ledger",
		Description: "Returns the ledger of buddy points transactions for the currently logged in user",
		Tags:        []string{"User"},
		Security:    []map[string][]string{{"session-cookie": {}}},
	}, h.getBuddyPointsLedger)
}

func (h *Handler) getBuddyPoints(ctx context.Context, input *struct{}) (*GetBuddyPointsOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized(http.StatusText(http.StatusUnauthorized))
	}

	points, err := h.service.GetBuddyPoints(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get buddy points", err)
	}

	return &GetBuddyPointsOutput{
		Body: GetBuddyPointsOutputBody{
			MonthlyPoints:  points.MonthlyPoints,
			LifetimePoints: points.LifetimePoints,
		},
	}, nil
}

func (h *Handler) getBuddyPointsLedger(ctx context.Context, input *struct{}) (*GetBuddyPointsLedgerOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized(http.StatusText(http.StatusUnauthorized))
	}

	ledger, err := h.service.GetBuddyPointsLedger(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get buddy points ledger", err)
	}

	return &GetBuddyPointsLedgerOutput{
		Body: GetBuddyPointsLedgerOutputBody{
			Entries: ledger.Entries,
		},
	}, nil
}
