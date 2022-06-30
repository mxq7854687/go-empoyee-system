// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: employee.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createEmployee = `-- name: CreateEmployee :one
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
    RETURNING employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id
`

type CreateEmployeeParams struct {
	FirstName    sql.NullString `json:"first_name"`
	LastName     string         `json:"last_name"`
	Email        string         `json:"email"`
	PhoneNumber  sql.NullString `json:"phone_number"`
	HireDate     time.Time      `json:"hire_date"`
	JobID        int64          `json:"job_id"`
	Salary       int64          `json:"salary"`
	ManagerID    sql.NullInt64  `json:"manager_id"`
	DepartmentID int64          `json:"department_id"`
}

func (q *Queries) CreateEmployee(ctx context.Context, arg CreateEmployeeParams) (Employee, error) {
	row := q.db.QueryRowContext(ctx, createEmployee,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.PhoneNumber,
		arg.HireDate,
		arg.JobID,
		arg.Salary,
		arg.ManagerID,
		arg.DepartmentID,
	)
	var i Employee
	err := row.Scan(
		&i.EmployeeID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNumber,
		&i.HireDate,
		&i.JobID,
		&i.Salary,
		&i.ManagerID,
		&i.DepartmentID,
	)
	return i, err
}

const deleteEmployee = `-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE employee_id = $1
`

func (q *Queries) DeleteEmployee(ctx context.Context, employeeID int64) error {
	_, err := q.db.ExecContext(ctx, deleteEmployee, employeeID)
	return err
}

const getEmployee = `-- name: GetEmployee :one
SELECT employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id FROM employees
WHERE employee_id = $1 LIMIT 1
`

func (q *Queries) GetEmployee(ctx context.Context, employeeID int64) (Employee, error) {
	row := q.db.QueryRowContext(ctx, getEmployee, employeeID)
	var i Employee
	err := row.Scan(
		&i.EmployeeID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNumber,
		&i.HireDate,
		&i.JobID,
		&i.Salary,
		&i.ManagerID,
		&i.DepartmentID,
	)
	return i, err
}

const listEmployees = `-- name: ListEmployees :many
SELECT employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id FROM employees
ORDER BY employee_id
LIMIT $1
OFFSET $2
`

type ListEmployeesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListEmployees(ctx context.Context, arg ListEmployeesParams) ([]Employee, error) {
	rows, err := q.db.QueryContext(ctx, listEmployees, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Employee{}
	for rows.Next() {
		var i Employee
		if err := rows.Scan(
			&i.EmployeeID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNumber,
			&i.HireDate,
			&i.JobID,
			&i.Salary,
			&i.ManagerID,
			&i.DepartmentID,
		); err != nil {
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

const updateEmployee = `-- name: UpdateEmployee :one
UPDATE employees
SET first_name = $2, last_name = $3, phone_number = $4,
    hire_date = $5, manager_id = $6, department_id = $7
WHERE employee_id = $1
RETURNING employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id
`

type UpdateEmployeeParams struct {
	EmployeeID   int64          `json:"employee_id"`
	FirstName    sql.NullString `json:"first_name"`
	LastName     string         `json:"last_name"`
	PhoneNumber  sql.NullString `json:"phone_number"`
	HireDate     time.Time      `json:"hire_date"`
	ManagerID    sql.NullInt64  `json:"manager_id"`
	DepartmentID int64          `json:"department_id"`
}

func (q *Queries) UpdateEmployee(ctx context.Context, arg UpdateEmployeeParams) (Employee, error) {
	row := q.db.QueryRowContext(ctx, updateEmployee,
		arg.EmployeeID,
		arg.FirstName,
		arg.LastName,
		arg.PhoneNumber,
		arg.HireDate,
		arg.ManagerID,
		arg.DepartmentID,
	)
	var i Employee
	err := row.Scan(
		&i.EmployeeID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PhoneNumber,
		&i.HireDate,
		&i.JobID,
		&i.Salary,
		&i.ManagerID,
		&i.DepartmentID,
	)
	return i, err
}
