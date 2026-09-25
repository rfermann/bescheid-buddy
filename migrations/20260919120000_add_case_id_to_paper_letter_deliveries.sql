-- +goose Up
-- +goose StatementBegin
ALTER TABLE paper_letter_deliveries
    ADD COLUMN case_id UUID REFERENCES cases (id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE paper_letter_deliveries
    DROP COLUMN case_id;
-- +goose StatementEnd
