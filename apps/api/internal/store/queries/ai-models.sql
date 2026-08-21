-- name: ListEnabledAIModels :many
SELECT id, provider_model_id, name, provider, enabled, sort_order, created_at, updated_at
FROM ai_models
WHERE enabled = TRUE
ORDER BY sort_order, created_at, id;

-- name: ListAIModels :many
SELECT id, provider_model_id, name, provider, enabled, sort_order, created_at, updated_at
FROM ai_models
ORDER BY sort_order, created_at, id;

-- name: GetAIModelByProviderModelID :one
SELECT id, provider_model_id, name, provider, enabled, sort_order, created_at, updated_at
FROM ai_models
WHERE provider_model_id = $1;

-- name: CreateAIModel :one
INSERT INTO ai_models (provider_model_id, name, provider, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING id, provider_model_id, name, provider, enabled, sort_order, created_at, updated_at;

-- name: SetAIModelEnabled :one
UPDATE ai_models
SET enabled    = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, provider_model_id, name, provider, enabled, sort_order, created_at, updated_at;
