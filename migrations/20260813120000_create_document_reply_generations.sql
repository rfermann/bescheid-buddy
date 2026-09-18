-- +goose Up
-- +goose StatementBegin
CREATE TYPE document_reply_generation_status AS enum ('pending', 'processing', 'completed', 'failed');

CREATE TYPE document_reply_generation_mode AS enum ('generate', 'improve');

CREATE TABLE document_reply_generations
(
    id                     uuid                              DEFAULT uuidv4() PRIMARY KEY,
    user_id                uuid                     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    document_id            uuid                              REFERENCES case_documents (id) ON DELETE SET NULL,
    mode                   document_reply_generation_mode   NOT NULL,
    client_request_id      uuid                              NOT NULL,
    status                 document_reply_generation_status NOT NULL DEFAULT 'pending',
    job_id                 bigint,
    model                  text                              NOT NULL,
    provider               text                              NOT NULL,
    language               text                              NOT NULL,
    input                  jsonb                             NOT NULL,
    result                 jsonb,
    error_message          text,
    started_at             timestamptz,
    completed_at           timestamptz,
    created_at             timestamptz DEFAULT NOW() NOT NULL,
    updated_at             timestamptz DEFAULT NOW() NOT NULL,
    CONSTRAINT document_reply_generations_completed_has_result
        CHECK (status <> 'completed' OR result IS NOT NULL),
    CONSTRAINT document_reply_generations_user_client_request_unique
        UNIQUE (user_id, client_request_id)
);

CREATE INDEX idx_document_reply_generations_user_created_at ON document_reply_generations (user_id, created_at);
CREATE INDEX idx_document_reply_generations_document_id ON document_reply_generations (document_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_document_reply_generations_user_created_at;
DROP INDEX IF EXISTS idx_document_reply_generations_document_id;
DROP TABLE IF EXISTS document_reply_generations;
DROP TYPE IF EXISTS document_reply_generation_mode;
DROP TYPE IF EXISTS document_reply_generation_status;
-- +goose StatementEnd
