-- name: CreateJob :one
INSERT INTO jobs (
    job_title,
    min_salary,
    max_salary
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE job_id = $1 LIMIT 1;

-- name: UpdateJob :one
UPDATE jobs
SET job_title = $2, min_salary = $3, max_salary = $4
WHERE job_id = $1
RETURNING *;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY job_id
LIMIT $1
OFFSET $2;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE job_id = $1;