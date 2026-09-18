-- +goose Up
-- +goose StatementBegin
CREATE TYPE paper_letter_upload_status AS enum ('pending', 'confirmed');

CREATE TABLE paper_letter_uploads
(
    id                 uuid                       DEFAULT uuidv4() PRIMARY KEY,
    user_id            uuid              NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    sanitized_filename varchar(255)      NOT NULL,
    status             paper_letter_upload_status NOT NULL DEFAULT 'pending',
    size_bytes         bigint            NOT NULL DEFAULT 0,
    page_count         integer           NOT NULL DEFAULT 0,
    sha256             text              NOT NULL DEFAULT '',
    confirmed_at       timestamptz,
    created_at          timestamptz        NOT NULL DEFAULT NOW(),
    updated_at          timestamptz        NOT NULL DEFAULT NOW(),
    CONSTRAINT paper_letter_uploads_nonnegative_metadata
        CHECK (size_bytes >= 0 AND page_count >= 0),
    CONSTRAINT paper_letter_uploads_sha256_format
        CHECK (sha256 = '' OR sha256 ~ '^[0-9a-f]{64}$'),
    CONSTRAINT paper_letter_uploads_confirmed_has_metadata
        CHECK (status <> 'confirmed' OR
               (size_bytes > 0 AND page_count > 0 AND sha256 <> '' AND confirmed_at IS NOT NULL))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS paper_letter_uploads;
DROP TYPE IF EXISTS paper_letter_upload_status;
-- +goose StatementEnd
