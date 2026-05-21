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
