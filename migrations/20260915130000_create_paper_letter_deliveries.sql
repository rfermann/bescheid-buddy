-- +goose Up
-- +goose StatementBegin
CREATE TYPE paper_letter_delivery_status AS enum ('pending', 'processing', 'sent', 'failed', 'manual_review');

CREATE TABLE paper_letter_deliveries
(
    id                             uuid                         DEFAULT uuidv4() PRIMARY KEY,
    user_id                        uuid                NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    upload_id                      uuid                NOT NULL REFERENCES paper_letter_uploads (id) ON DELETE RESTRICT,
    client_request_id              uuid                NOT NULL,
    status                         paper_letter_delivery_status NOT NULL DEFAULT 'pending',
    job_id                         bigint,
    pingen_id                      text,
    provider_submission_started_at timestamptz,
    page_count                     integer             NOT NULL,
    bp_cost                        integer             NOT NULL,
    delivery_product               text                NOT NULL,
    address_position               text                NOT NULL,
    print_mode                     text                NOT NULL,
    print_spectrum                 text                NOT NULL,
    error_message                  text,
    sent_at                        timestamptz,
    created_at                     timestamptz         NOT NULL DEFAULT NOW(),
    updated_at                     timestamptz         NOT NULL DEFAULT NOW(),
    CONSTRAINT paper_letter_deliveries_positive_values CHECK (page_count > 0 AND bp_cost > 0),
    CONSTRAINT paper_letter_deliveries_sent_has_provider_id
        CHECK (status <> 'sent' OR (pingen_id IS NOT NULL AND sent_at IS NOT NULL)),
    CONSTRAINT paper_letter_deliveries_user_client_request_unique UNIQUE (user_id, client_request_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS paper_letter_deliveries;
DROP TYPE IF EXISTS paper_letter_delivery_status;
-- +goose StatementEnd
