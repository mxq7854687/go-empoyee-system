
-- name: CreateRole :one
INSERT INTO roles (
role,
privileges
) VALUES (
$1, $2
)
RETURNING *;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;

-- name: GetRoleByRoleName :one
SELECT * FROM roles
WHERE role = $1 LIMIT 1;

-- name: DeleteAllRole :exec
DELETE FROM roles;