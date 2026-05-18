package buddy_package_handler

import (
	"context"
	"net/http"

	"github.com/bored-bots/bored-template/internal/services/buddy_package"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *buddy_package.Service
}

func NewHandler(service *buddy_package.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type GetAllBuddyPackagesOutputBody struct {
	Packages []*db.BuddyPackage `json:"packages"`
}

type GetAllBuddyPackagesOutput struct {
	Body GetAllBuddyPackagesOutputBody `json:"body"`
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-all-buddy-packages",
		Method:      http.MethodGet,
		Path:        "/v1/buddy-packages",
		Summary:     "Get all buddy packages",
		Description: "Returns a list of all buddy packages sorted by price ascending",
		Tags:        []string{"BuddyPackage"},
	}, h.getAllBuddyPackages)
}

func (h *Handler) getAllBuddyPackages(ctx context.Context, input *struct{}) (*GetAllBuddyPackagesOutput, error) {
	packages, err := h.service.GetAllBuddyPackages(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get buddy packages", err)
	}

	return &GetAllBuddyPackagesOutput{
		Body: GetAllBuddyPackagesOutputBody{
			Packages: packages,
		},
	}, nil
}
