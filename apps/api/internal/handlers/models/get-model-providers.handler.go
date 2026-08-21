package models_handler

import "context"

type GetModelProvidersOutputBody struct {
	Providers []string `json:"providers"`
}

type GetModelProvidersOutput struct {
	Body GetModelProvidersOutputBody `json:"body"`
}

func (h *Handler) getModelProviders(_ context.Context, _ *struct{}) (*GetModelProvidersOutput, error) {
	return &GetModelProvidersOutput{
		Body: GetModelProvidersOutputBody{Providers: h.service.GetProviders()},
	}, nil
}
