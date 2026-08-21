-- +goose Up
-- +goose StatementBegin
CREATE TABLE ai_models
(
    id                uuid        DEFAULT uuidv4() PRIMARY KEY,
    provider_model_id text                      NOT NULL UNIQUE,
    name              text                      NOT NULL,
    provider          text                      NOT NULL,
    enabled           boolean     DEFAULT TRUE  NOT NULL,
    sort_order        integer     DEFAULT 0     NOT NULL,
    created_at        timestamptz DEFAULT NOW() NOT NULL,
    updated_at        timestamptz DEFAULT NOW() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ai_models;
-- +goose StatementEnd
