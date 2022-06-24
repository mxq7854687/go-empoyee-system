
-- name: CreateUser :one
INSERT INTO users (
email,
hashed_password,
status
) VALUES (
$1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;