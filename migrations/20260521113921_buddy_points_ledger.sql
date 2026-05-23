-- +goose Up
CREATE TYPE transaction_type AS enum (
    'purchase',
    'referral',
    'monthly_reset',
    'usage'
    );

CREATE TYPE related_entity_type AS enum (
    'buddy_package',
    'user',
    'null'
    );

CREATE TABLE IF NOT EXISTS buddy_points_ledger
(
    id                     uuid PRIMARY KEY DEFAULT uuidv4() NOT NULL,
    user_id                uuid                              NOT NULL,
    transaction_type       transaction_type                  NOT NULL,
    amount                 integer                           NOT NULL,
    monthly_points_before  integer                           NOT NULL,
    lifetime_points_before integer                           NOT NULL,
    monthly_points_after   integer                           NOT NULL,
    lifetime_points_after  integer                           NOT NULL,
    description            text                              NOT NULL,
    related_entity_id      uuid,
    related_entity_type    related_entity_type,
    created_at             timestamptz      DEFAULT NOW()    NOT NULL,
    CONSTRAINT buddy_points_ledger_user_id_fk FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE INDEX idx_buddy_points_ledger_user_id ON buddy_points_ledger (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_buddy_points_ledger_user_id;
DROP TABLE IF EXISTS buddy_points_ledger CASCADE;
DROP TYPE IF EXISTS related_entity_type;
DROP TYPE IF EXISTS transaction_type;