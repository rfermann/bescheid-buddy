package user_handler

import (
	"context"
	"fmt"
	"net/http"

	appContext "github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/services/user_service"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *user_service.Service
}

func NewHandler(service *user_service.Service) *Handler {
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
}

func (h *Handler) getBuddyPoints(ctx context.Context, input *struct{}) (*GetBuddyPointsOutput, error) {
	userID, ok := appContext.UserIDFromContext(ctx)
	fmt.Println("userID", userID)
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
