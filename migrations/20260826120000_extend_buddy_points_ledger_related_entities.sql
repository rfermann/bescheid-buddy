-- +goose Up
-- +goose StatementBegin
ALTER TYPE related_entity_type ADD VALUE 'paper_letter_generation';
ALTER TYPE related_entity_type ADD VALUE 'paper_letter_delivery';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM buddy_points_ledgers
        WHERE related_entity_type::text IN ('paper_letter_generation', 'paper_letter_delivery')
    ) THEN
        RAISE EXCEPTION 'cannot remove paper-letter related entity types while ledger rows use them';
    END IF;
END
$$;

ALTER TYPE related_entity_type RENAME TO related_entity_type_with_paper_letters;

CREATE TYPE related_entity_type AS enum (
    'buddy_packages',
    'user',
    'null'
);

ALTER TABLE buddy_points_ledgers
    ALTER COLUMN related_entity_type TYPE related_entity_type
        USING (related_entity_type::text::related_entity_type);

DROP TYPE related_entity_type_with_paper_letters;
-- +goose StatementEnd
