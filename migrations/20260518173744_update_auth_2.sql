-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user" ADD COLUMN "stripe_customer_id" text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user" DROP COLUMN "stripe_customer_id";
-- +goose StatementEnd
