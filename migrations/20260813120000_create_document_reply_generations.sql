-- +goose Up
-- +goose StatementBegin
CREATE TYPE document_reply_generation_status AS enum ('pending', 'processing', 'completed', 'failed');

CREATE TABLE document_reply_generations
(
    id                     uuid                             DEFAULT uuidv4() PRIMARY KEY,
    document_id            uuid                            NOT NULL REFERENCES case_documents (id) ON DELETE CASCADE,
    status                 document_reply_generation_status NOT NULL DEFAULT 'pending',
    job_id                 bigint,
    model                  text                      NOT NULL,
    provider               text                      NOT NULL,
    language               text                      NOT NULL,
    instruction            text                      NOT NULL DEFAULT '',
    result                 jsonb,
    error_message          text,
    started_at             timestamptz,
    completed_at           timestamptz,
    created_at             timestamptz DEFAULT NOW() NOT NULL,
    updated_at             timestamptz DEFAULT NOW() NOT NULL,
    CONSTRAINT document_reply_generations_completed_has_result
        CHECK (status <> 'completed' OR result IS NOT NULL)
);

CREATE INDEX idx_document_reply_generations_document_id ON document_reply_generations (document_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_document_reply_generations_document_id;
DROP TABLE IF EXISTS document_reply_generations;
DROP TYPE IF EXISTS document_reply_generation_status;
-- +goose StatementEnd
