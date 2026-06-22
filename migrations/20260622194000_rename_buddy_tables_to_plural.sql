-- +goose Up
-- +goose StatementBegin
-- Rename enum value
ALTER TYPE related_entity_type RENAME VALUE 'buddy_package' TO 'buddy_packages';

-- Drop foreign key constraint before renaming tables
ALTER TABLE buddy_points_ledger
    DROP CONSTRAINT buddy_points_ledger_user_id_fk;

-- Rename tables
ALTER TABLE buddy_package
    RENAME TO buddy_packages;
ALTER TABLE buddy_points_ledger
    RENAME TO buddy_points_ledgers;

-- Rename index
ALTER INDEX idx_buddy_points_ledger_user_id RENAME TO idx_buddy_points_ledgers_user_id;

-- Re-add foreign key constraint with new name
ALTER TABLE buddy_points_ledgers
    ADD CONSTRAINT buddy_points_ledgers_user_id_fk
        FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE;

-- Rename unique constraint
ALTER TABLE buddy_packages
    RENAME CONSTRAINT buddy_package_stripe_price_id_unique TO buddy_packages_stripe_price_id_unique;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Reverse the renames
ALTER TABLE buddy_packages
    RENAME CONSTRAINT buddy_packages_stripe_price_id_unique TO buddy_package_stripe_price_id_unique;

ALTER TABLE buddy_points_ledgers
    DROP CONSTRAINT buddy_points_ledgers_user_id_fk;

ALTER TABLE buddy_packages
    RENAME TO buddy_package;
ALTER TABLE buddy_points_ledgers
    RENAME TO buddy_points_ledger;

ALTER INDEX idx_buddy_points_ledgers_user_id RENAME TO idx_buddy_points_ledger_user_id;

ALTER TABLE buddy_points_ledger
    ADD CONSTRAINT buddy_points_ledger_user_id_fk
        FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE;

ALTER TYPE related_entity_type RENAME VALUE 'buddy_packages' TO 'buddy_package';
-- +goose StatementEnd
