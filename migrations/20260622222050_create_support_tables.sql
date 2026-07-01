-- +goose Up
-- +goose StatementBegin
CREATE TYPE support_case_status AS enum ('open', 'in_progress', 'resolved', 'closed');
CREATE TYPE support_case_category AS enum ('general', 'technical', 'billing', 'feature_request');

CREATE TABLE support_cases
(
    id          uuid        DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id     uuid                      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    subject     text                      NOT NULL,
    status      support_case_status       NOT NULL,
    category    support_case_category     NOT NULL,
    updated_by  uuid                      NOT NULL REFERENCES users (id) ON DELETE SET NULL,
    created_at  timestamptz DEFAULT NOW() NOT NULL,
    updated_at  timestamptz DEFAULT NOW() NOT NULL,
    resolved_at timestamptz
);

CREATE TABLE support_case_messages
(
    id              uuid        DEFAULT gen_random_uuid() PRIMARY KEY,
    case_id         uuid                      NOT NULL REFERENCES support_cases (id) ON DELETE CASCADE,
    user_id         uuid                      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    message         text                      NOT NULL,
    is_user_message boolean     DEFAULT FALSE NOT NULL,
    is_internal     boolean     DEFAULT FALSE NOT NULL,
    is_read         boolean     DEFAULT FALSE NOT NULL,
    updated_by      uuid                      NOT NULL REFERENCES users (id) ON DELETE SET NULL,
    created_at      timestamptz DEFAULT NOW() NOT NULL,
    updated_at      timestamptz DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_support_cases_user_id ON support_cases (user_id);
CREATE INDEX idx_support_messages_case_id ON support_case_messages (case_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_support_cases_user_id;
DROP INDEX IF EXISTS idx_support_messages_case_id;

DROP TABLE IF EXISTS support_case_messages;
DROP TABLE IF EXISTS support_cases;

DROP TYPE IF EXISTS support_case_status;
DROP TYPE IF EXISTS support_case_category;
-- +goose StatementEnd
