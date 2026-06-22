-- +goose Up
-- +goose StatementBegin
ALTER TABLE "account" RENAME TO "accounts";

ALTER TABLE "apikey" RENAME TO "apikeys";

ALTER TABLE "session" RENAME TO "sessions";

ALTER TABLE "user" RENAME TO "users";

ALTER TABLE "verification" RENAME TO "verifications";

ALTER TABLE "sessions" DROP CONSTRAINT "session_token_unique";

ALTER TABLE "users" DROP CONSTRAINT "user_email_unique";

ALTER TABLE "accounts" DROP CONSTRAINT "account_user_id_user_id_fk";


ALTER TABLE "sessions" DROP CONSTRAINT "session_user_id_user_id_fk";


DROP INDEX "account_userId_idx";

DROP INDEX "apikey_configId_idx";

DROP INDEX "apikey_referenceId_idx";

DROP INDEX "apikey_key_idx";

DROP INDEX "session_userId_idx";

DROP INDEX "verification_identifier_idx";

ALTER TABLE "accounts" ADD CONSTRAINT "accounts_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;

ALTER TABLE "sessions" ADD CONSTRAINT "sessions_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;

CREATE INDEX "accounts_userId_idx" ON "accounts" USING btree ("user_id");

CREATE INDEX "apikeys_configId_idx" ON "apikeys" USING btree ("config_id");

CREATE INDEX "apikeys_referenceId_idx" ON "apikeys" USING btree ("reference_id");

CREATE INDEX "apikeys_key_idx" ON "apikeys" USING btree ("key");

CREATE INDEX "sessions_userId_idx" ON "sessions" USING btree ("user_id");

CREATE INDEX "verifications_identifier_idx" ON "verifications" USING btree ("identifier");

ALTER TABLE "sessions" ADD CONSTRAINT "sessions_token_unique" UNIQUE("token");

ALTER TABLE "users" ADD CONSTRAINT "users_email_unique" UNIQUE("email");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users" DROP CONSTRAINT "users_email_unique";
ALTER TABLE "sessions" DROP CONSTRAINT "sessions_token_unique";
DROP INDEX IF EXISTS "verifications_identifier_idx";
DROP INDEX IF EXISTS "sessions_userId_idx";
DROP INDEX IF EXISTS "apikeys_key_idx";
DROP INDEX IF EXISTS "apikeys_referenceId_idx";
DROP INDEX IF EXISTS "apikeys_configId_idx";
DROP INDEX IF EXISTS "accounts_userId_idx";
ALTER TABLE "sessions" DROP CONSTRAINT "sessions_user_id_users_id_fk";
ALTER TABLE "accounts" DROP CONSTRAINT "accounts_user_id_users_id_fk";
CREATE INDEX "verification_identifier_idx";
CREATE INDEX "session_userId_idx";
CREATE INDEX "apikey_key_idx";
CREATE INDEX "apikey_referenceId_idx";
CREATE INDEX "apikey_configId_idx";
CREATE INDEX "account_userId_idx";
ALTER TABLE "sessions" ADD CONSTRAINT "session_user_id_user_id_fk";
ALTER TABLE "accounts" ADD CONSTRAINT "account_user_id_user_id_fk";
ALTER TABLE "users" ADD CONSTRAINT "user_email_unique";
ALTER TABLE "sessions" ADD CONSTRAINT "session_token_unique";
ALTER TABLE "verifications" RENAME TO "verification";
ALTER TABLE "users" RENAME TO "user";
ALTER TABLE "sessions" RENAME TO "session";
ALTER TABLE "apikeys" RENAME TO "apikey";
ALTER TABLE "accounts" RENAME TO "account";
-- +goose StatementEnd
