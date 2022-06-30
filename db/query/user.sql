
-- name: CreateUser :one
INSERT INTO users (
email,
hashed_password,
status,
role_id
) VALUES (
$1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ActivateUser :exec
UPDATE users
SET hashed_password = $1, status = 'activated', updated_at = now()
WHERE email = $2;