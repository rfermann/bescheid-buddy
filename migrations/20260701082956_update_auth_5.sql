-- +goose Up
-- +goose StatementBegin
ALTER TABLE "sessions" ADD COLUMN "impersonated_by" text;

ALTER TABLE "users" ADD COLUMN "role" text;

ALTER TABLE "users" ADD COLUMN "banned" boolean DEFAULT false;

ALTER TABLE "users" ADD COLUMN "ban_reason" text;

ALTER TABLE "users" ADD COLUMN "ban_expires" timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users" DROP COLUMN "ban_expires";
ALTER TABLE "users" DROP COLUMN "ban_reason";
ALTER TABLE "users" DROP COLUMN "banned";
ALTER TABLE "users" DROP COLUMN "role";
ALTER TABLE "sessions" DROP COLUMN "impersonated_by";
-- +goose StatementEnd
