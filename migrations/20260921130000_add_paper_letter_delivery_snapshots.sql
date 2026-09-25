-- +goose Up
-- +goose StatementBegin
ALTER TABLE paper_letter_deliveries
    ADD COLUMN recipient varchar(200),
    ADD COLUMN topic varchar(500);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE paper_letter_deliveries
    DROP COLUMN topic,
    DROP COLUMN recipient;
-- +goose StatementEnd
