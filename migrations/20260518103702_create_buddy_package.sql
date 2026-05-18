-- +goose Up
CREATE TABLE IF NOT EXISTS buddy_package
(
    id              uuid PRIMARY KEY DEFAULT uuidv4() NOT NULL,
    name            text                              NOT NULL,
    points          integer                           NOT NULL,
    price_in_cents  integer                           NOT NULL,
    stripe_price_id text                              NOT NULL,
    created_at      timestamptz      DEFAULT NOW()    NOT NULL,
    updated_at      timestamptz      DEFAULT NOW()    NOT NULL,
    CONSTRAINT buddy_package_stripe_price_id_unique UNIQUE (stripe_price_id)
);

-- +goose Down
DROP TABLE IF EXISTS buddy_package CASCADE;

