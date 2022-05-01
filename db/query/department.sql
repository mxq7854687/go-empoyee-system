-- name: CreateDepartment :one
INSERT INTO departments (
    department_name
) VALUES (
             $1
         )
    RETURNING *;

-- name: GetDepartment :one
SELECT * FROM departments
WHERE department_id = $1 LIMIT 1;

-- name: UpdateDepartment :one
UPDATE departments
set department_name = $1
WHERE department_id = $2
RETURNING *;

-- name: ListDepartments :many
SELECT * FROM departments
ORDER BY department_id
LIMIT $1
OFFSET $2;

-- name: DeleteDepartments :exec
DELETE FROM departments
WHERE department_id = $1;