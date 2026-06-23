-- name: GetAllBuddyPackages :many
SELECT id, name, points, price_in_cents, stripe_price_id, created_at, updated_at
FROM buddy_packages
ORDER BY price_in_cents;

-- name: GetBuddyPackageByID :one
SELECT id, name, points, price_in_cents, stripe_price_id, created_at, updated_at
FROM buddy_packages
WHERE id = $1;
