-- name: CreateSupportCase :one
INSERT INTO support_cases (user_id, subject, status, category, updated_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateSupportCaseMessage :one
INSERT INTO support_case_messages (case_id, user_id, message, is_user_message, is_internal, updated_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSupportCaseByID :one
SELECT *
FROM support_cases
WHERE id = $1
  AND user_id = $2;

-- name: GetSupportCase :one
SELECT *
FROM support_cases
WHERE id = $1;

-- name: GetSupportCasesByUserID :many
SELECT *
FROM support_cases
WHERE user_id = $1;

-- name: GetSupportCaseMessagesBySupportCaseID :many
SELECT *
FROM support_case_messages AS smc
WHERE case_id = $1
  AND EXISTS(SELECT 1
             FROM support_cases
             WHERE id = smc.case_id
               AND support_cases.user_id = $2);

-- name: GetAllSupportCases :many
SELECT *
FROM support_cases;

-- name: UpdateSupportCaseStatus :exec
UPDATE support_cases
SET status = $2,
    updated_by = $3,
    updated_at = NOW()
WHERE id = $1;