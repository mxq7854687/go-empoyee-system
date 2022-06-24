-- name: CreateEmployee :one
INSERT INTO employees (
	first_name,
	last_name,
	email,
	phone_number,
	hire_date,
	job_id,
	salary,
	manager_id,
	department_id
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         )
    RETURNING *;

-- name: GetEmployee :one
SELECT * FROM employees
WHERE employee_id = $1 LIMIT 1;

-- name: UpdateEmployee :one
UPDATE employees
SET first_name = $2, last_name = $3, phone_number = $4,
    hire_date = $5, manager_id = $6, department_id = $7
WHERE employee_id = $1
RETURNING *;

-- name: ListEmployees :many
SELECT * FROM employees
ORDER BY employee_id
LIMIT $1
OFFSET $2;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE employee_id = $1;