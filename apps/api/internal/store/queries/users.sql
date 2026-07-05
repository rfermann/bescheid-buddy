-- name: UpdateUserPoints :exec
UPDATE "users"
SET lifetime_points = lifetime_points + $2,
    updated_at      = NOW()
WHERE id = $1;

-- name: UpdateUserStripeCustomerID :exec
UPDATE "users"
SET stripe_customer_id = $2,
    updated_at         = NOW()
WHERE id = $1;

-- name: GetUserByID :one
SELECT *
FROM "users"
WHERE id = $1;

-- name: GetUserByStripeCustomerID :one
SELECT *
FROM "users"
WHERE stripe_customer_id = $1;

-- name: GetBuddyPointsByUserID :one
SELECT monthly_points,
       lifetime_points
FROM "users"
WHERE id = $1;

-- name: InsertBuddyPointsLedger :exec
INSERT INTO buddy_points_ledgers (user_id,
                                 transaction_type,
                                 amount,
                                 monthly_points_before,
                                 lifetime_points_before,
                                 monthly_points_after,
                                 lifetime_points_after,
                                 description,
                                 related_entity_id,
                                 related_entity_type)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10);

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
FROM buddy_points_ledgers
WHERE user_id = $1
ORDER BY created_at DESC;
