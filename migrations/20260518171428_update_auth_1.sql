-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user" ADD COLUMN "monthly_points" integer DEFAULT 5 NOT NULL;

ALTER TABLE "user" ADD COLUMN "lifetime_points" integer DEFAULT 0 NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user" DROP COLUMN "lifetime_points";
ALTER TABLE "user" DROP COLUMN "monthly_points";
-- +goose StatementEnd
