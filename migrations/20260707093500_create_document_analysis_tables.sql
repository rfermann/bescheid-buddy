-- +goose Up
-- +goose StatementBegin
CREATE TYPE document_analysis_status AS enum ('pending', 'processing', 'completed', 'failed');

CREATE TYPE document_upload_status AS enum ('awaiting_upload', 'uploaded', 'failed');

CREATE TABLE cases
(
    id         uuid        DEFAULT uuidv4() PRIMARY KEY,
    user_id    uuid                      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      varchar(255),
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL
);

CREATE TABLE case_documents
(
    id            uuid                              DEFAULT uuidv4() PRIMARY KEY,
    case_id       uuid                     NOT NULL REFERENCES cases (id) ON DELETE CASCADE,
    file_name     varchar(255)             NOT NULL,
    upload_status document_upload_status   NOT NULL DEFAULT 'awaiting_upload',
    uploaded_at   timestamptz,
    status        document_analysis_status NOT NULL DEFAULT 'pending',
    error_message text,
    created_at    timestamptz                       DEFAULT NOW() NOT NULL,
    updated_at    timestamptz                       DEFAULT NOW() NOT NULL,
    completed_at  timestamptz
);

CREATE TABLE document_extractions
(
    id            uuid        DEFAULT uuidv4() PRIMARY KEY,
    document_id   uuid                      NOT NULL UNIQUE REFERENCES case_documents (id) ON DELETE CASCADE,
    status        document_analysis_status  NOT NULL DEFAULT 'pending',
    job_id        bigint,
    text          text,
    provider      text,
    mime          text,
    page_count    integer,
    char_count    integer,
    text_sha256   text,
    error_message text,
    started_at    timestamptz,
    completed_at  timestamptz,
    created_at    timestamptz DEFAULT NOW() NOT NULL,
    updated_at    timestamptz DEFAULT NOW() NOT NULL,
    CONSTRAINT document_extractions_completed_has_text
        CHECK (status <> 'completed' OR text IS NOT NULL)
);

CREATE TABLE document_analysis_results
(
    id                uuid        DEFAULT uuidv4() PRIMARY KEY,
    document_id       uuid                      NOT NULL UNIQUE REFERENCES case_documents (id) ON DELETE CASCADE,
    status            document_analysis_status  NOT NULL DEFAULT 'pending',
    job_id            bigint,
    model             text,
    language          text,
    detected_language text,
    provider          text,
    analysis_text     text,
    analysis_metadata jsonb,
    error_message     text,
    started_at        timestamptz,
    completed_at      timestamptz,
    created_at        timestamptz DEFAULT NOW() NOT NULL,
    updated_at        timestamptz DEFAULT NOW() NOT NULL,
    CONSTRAINT document_analysis_results_completed_has_text
        CHECK (status <> 'completed' OR analysis_text IS NOT NULL)
);

CREATE TABLE case_conversations
(
    id              uuid        DEFAULT uuidv4() PRIMARY KEY,
    case_id         uuid                      NOT NULL REFERENCES cases (id) ON DELETE CASCADE,
    document_id     uuid                      REFERENCES case_documents (id) ON DELETE SET NULL,
    question        text                      NOT NULL,
    answer          text                      NOT NULL,
    answer_metadata jsonb,
    created_at      timestamptz DEFAULT NOW() NOT NULL,
    updated_at      timestamptz DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_cases_user_id ON cases (user_id);
CREATE INDEX idx_case_documents_case_id ON case_documents (case_id);
CREATE INDEX idx_document_extractions_document_id ON document_extractions (document_id);
CREATE INDEX idx_document_analysis_results_document_id ON document_analysis_results (document_id);
CREATE INDEX idx_document_conversations_case_id ON case_conversations (case_id);
CREATE INDEX idx_document_conversations_document_id ON case_conversations (document_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cases_user_id;
DROP INDEX IF EXISTS idx_case_documents_case_id;
DROP INDEX IF EXISTS idx_document_extractions_document_id;
DROP INDEX IF EXISTS idx_document_analysis_results_document_id;
DROP INDEX IF EXISTS idx_document_conversations_case_id;
DROP INDEX IF EXISTS idx_document_conversations_document_id;

DROP TABLE IF EXISTS case_conversations;
DROP TABLE IF EXISTS document_analysis_results;
DROP TABLE IF EXISTS document_extractions;
DROP TABLE IF EXISTS case_documents;
DROP TABLE IF EXISTS cases;

DROP TYPE IF EXISTS document_analysis_status;
DROP TYPE IF EXISTS document_upload_status;
-- +goose StatementEnd
