-- name: UpdateUserPoints :exec
UPDATE "user"
SET lifetime_points = lifetime_points + $2,
    updated_at      = NOW()
WHERE id = $1;

-- name: GetUserByID :one
SELECT id,
       name,
       email,
       email_verified,
       image,
       created_at,
       updated_at,
       telegram_id,
       telegram_phone_number,
       telegram_username,
       monthly_points,
       lifetime_points,
       stripe_customer_id
FROM "user"
WHERE id = $1;

-- name: GetUserByStripeCustomerID :one
SELECT id,
       name,
       email,
       email_verified,
       image,
       created_at,
       updated_at,
       telegram_id,
       telegram_phone_number,
       telegram_username,
       monthly_points,
       lifetime_points,
       stripe_customer_id
FROM "user"
WHERE stripe_customer_id = $1;

-- name: GetBuddyPointsByUserID :one
SELECT monthly_points,
       lifetime_points
FROM "user"
WHERE id = $1;

-- name: InsertBuddyPointsLedger :exec
INSERT INTO buddy_points_ledger (
    user_id,
    transaction_type,
    amount,
    monthly_points_before,
    lifetime_points_before,
    monthly_points_after,
    lifetime_points_after,
    description,
    related_entity_id,
    related_entity_type
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
);

-- name: GetBuddyPointsLedgerByUserID :many
SELECT id,
       user_id,
       transaction_type,
       amount,
       monthly_points_before,
       lifetime_points_before,
       monthly_points_after,
       lifetime_points_after,
       description,
       related_entity_id,
       related_entity_type,
       created_at
FROM buddy_points_ledger
WHERE user_id = $1
ORDER BY created_at DESC;
