// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: department.sql

package db

import (
	"context"
)

const createDepartment = `-- name: CreateDepartment :one
INSERT INTO departments (
    department_name
) VALUES (
             $1
         )
    RETURNING department_id, department_name
`

func (q *Queries) CreateDepartment(ctx context.Context, departmentName string) (Department, error) {
	row := q.db.QueryRowContext(ctx, createDepartment, departmentName)
	var i Department
	err := row.Scan(&i.DepartmentID, &i.DepartmentName)
	return i, err
}

const deleteDepartments = `-- name: DeleteDepartments :exec
DELETE FROM departments
WHERE department_id = $1
`

func (q *Queries) DeleteDepartments(ctx context.Context, departmentID int64) error {
	_, err := q.db.ExecContext(ctx, deleteDepartments, departmentID)
	return err
}

const getDepartment = `-- name: GetDepartment :one
SELECT department_id, department_name FROM departments
WHERE department_id = $1 LIMIT 1
`

func (q *Queries) GetDepartment(ctx context.Context, departmentID int64) (Department, error) {
	row := q.db.QueryRowContext(ctx, getDepartment, departmentID)
	var i Department
	err := row.Scan(&i.DepartmentID, &i.DepartmentName)
	return i, err
}

const listDepartments = `-- name: ListDepartments :many
SELECT department_id, department_name FROM departments
ORDER BY department_id
LIMIT $1
OFFSET $2
`

type ListDepartmentsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListDepartments(ctx context.Context, arg ListDepartmentsParams) ([]Department, error) {
	rows, err := q.db.QueryContext(ctx, listDepartments, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Department
	for rows.Next() {
		var i Department
		if err := rows.Scan(&i.DepartmentID, &i.DepartmentName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateDepartment = `-- name: UpdateDepartment :one
UPDATE departments
set department_name = $1
WHERE department_id = $2
RETURNING department_id, department_name
`

type UpdateDepartmentParams struct {
	DepartmentName string `json:"department_name"`
	DepartmentID   int64  `json:"department_id"`
}

func (q *Queries) UpdateDepartment(ctx context.Context, arg UpdateDepartmentParams) (Department, error) {
	row := q.db.QueryRowContext(ctx, updateDepartment, arg.DepartmentName, arg.DepartmentID)
	var i Department
	err := row.Scan(&i.DepartmentID, &i.DepartmentName)
	return i, err
}