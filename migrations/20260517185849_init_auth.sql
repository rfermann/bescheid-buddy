-- +goose Up
-- +goose StatementBegin
CREATE TABLE "account"
(
    "id"                       uuid PRIMARY KEY DEFAULT pg_catalog.gen_random_uuid() NOT NULL,
    "account_id"               text                                                  NOT NULL,
    "provider_id"              text                                                  NOT NULL,
    "user_id"                  uuid                                                  NOT NULL,
    "access_token"             text,
    "refresh_token"            text,
    "id_token"                 text,
    "access_token_expires_at"  timestamp,
    "refresh_token_expires_at" timestamp,
    "scope"                    text,
    "password"                 text,
    "created_at"               timestamp        DEFAULT NOW()                        NOT NULL,
    "updated_at"               timestamp                                             NOT NULL,
    "telegram_id"              text,
    "telegram_username"        text
);


CREATE TABLE "session"
(
    "id"         uuid PRIMARY KEY DEFAULT pg_catalog.gen_random_uuid() NOT NULL,
    "expires_at" timestamp                                             NOT NULL,
    "token"      text                                                  NOT NULL,
    "created_at" timestamp        DEFAULT NOW()                        NOT NULL,
    "updated_at" timestamp                                             NOT NULL,
    "ip_address" text,
    "user_agent" text,
    "user_id"    uuid                                                  NOT NULL,
    CONSTRAINT "session_token_unique" UNIQUE ("token")
);


CREATE TABLE "user"
(
    "id"                    uuid PRIMARY KEY DEFAULT pg_catalog.gen_random_uuid() NOT NULL,
    "name"                  text                                                  NOT NULL,
    "email"                 text                                                  NOT NULL,
    "email_verified"        boolean          DEFAULT FALSE                        NOT NULL,
    "image"                 text,
    "created_at"            timestamp        DEFAULT NOW()                        NOT NULL,
    "updated_at"            timestamp        DEFAULT NOW()                        NOT NULL,
    "telegram_id"           text,
    "telegram_phone_number" text,
    "telegram_username"     text,
    CONSTRAINT "user_email_unique" UNIQUE ("email")
);


CREATE TABLE "verification"
(
    "id"         uuid PRIMARY KEY DEFAULT pg_catalog.gen_random_uuid() NOT NULL,
    "identifier" text                                                  NOT NULL,
    "value"      text                                                  NOT NULL,
    "expires_at" timestamp                                             NOT NULL,
    "created_at" timestamp        DEFAULT NOW()                        NOT NULL,
    "updated_at" timestamp        DEFAULT NOW()                        NOT NULL
);


ALTER TABLE "account"
    ADD CONSTRAINT "account_user_id_user_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."user" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "session"
    ADD CONSTRAINT "session_user_id_user_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."user" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

CREATE INDEX "account_userId_idx" ON "account" USING btree ("user_id");

CREATE INDEX "session_userId_idx" ON "session" USING btree ("user_id");

CREATE INDEX "verification_identifier_idx" ON "verification" USING btree ("identifier");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "verification_identifier_idx";
DROP INDEX IF EXISTS "session_userId_idx";
DROP INDEX IF EXISTS "account_userId_idx";
ALTER TABLE "session"
    DROP CONSTRAINT "session_user_id_user_id_fk";
ALTER TABLE "account"
    DROP CONSTRAINT "account_user_id_user_id_fk";
DROP TABLE IF EXISTS "verification" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "session" CASCADE;
DROP TABLE IF EXISTS "account" CASCADE;
-- +goose StatementEnd
