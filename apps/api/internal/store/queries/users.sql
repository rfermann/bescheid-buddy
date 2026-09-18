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

-- name: GetBuddyPointsForUpdate :one
SELECT monthly_points,
       lifetime_points
FROM "users"
WHERE id = sqlc.arg(user_id)
FOR UPDATE;

-- name: DebitBuddyPoints :execrows
UPDATE "users"
SET monthly_points  = monthly_points - sqlc.arg(monthly_points),
    lifetime_points = lifetime_points - sqlc.arg(lifetime_points),
    updated_at      = NOW()
WHERE id = sqlc.arg(user_id)
  AND monthly_points >= sqlc.arg(monthly_points)
  AND lifetime_points >= sqlc.arg(lifetime_points);

-- name: GetMatchingBuddyPointSpend :one
SELECT amount,
       monthly_points_before,
       lifetime_points_before,
       monthly_points_after,
       lifetime_points_after
FROM buddy_points_ledgers
WHERE user_id = sqlc.arg(user_id)
  AND transaction_type = 'usage'
  AND amount < 0
  AND related_entity_id = sqlc.arg(related_entity_id)
  AND related_entity_type = sqlc.arg(related_entity_type)
ORDER BY created_at ASC
LIMIT 1;

-- name: CreditBuddyPoints :execrows
UPDATE "users"
SET monthly_points  = monthly_points + sqlc.arg(monthly_points),
    lifetime_points = lifetime_points + sqlc.arg(lifetime_points),
    updated_at      = NOW()
WHERE id = sqlc.arg(user_id)
  AND sqlc.arg(monthly_points) >= 0
  AND sqlc.arg(lifetime_points) >= 0;

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
