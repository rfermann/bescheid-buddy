-- name: GetAllBuddyPackages :many
SELECT id, name, points, price_in_cents, stripe_price_id, created_at, updated_at
FROM buddy_package
ORDER BY price_in_cents;
