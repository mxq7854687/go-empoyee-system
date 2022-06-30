// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"context"
)

type Querier interface {
	ActivateUser(ctx context.Context, arg ActivateUserParams) error
	CreateDepartment(ctx context.Context, departmentName string) (Department, error)
	CreateEmployee(ctx context.Context, arg CreateEmployeeParams) (Employee, error)
	CreateJob(ctx context.Context, arg CreateJobParams) (Job, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteDepartments(ctx context.Context, departmentID int64) error
	DeleteEmployee(ctx context.Context, employeeID int64) error
	DeleteJob(ctx context.Context, jobID int64) error
	GetDepartment(ctx context.Context, departmentID int64) (Department, error)
	GetEmployee(ctx context.Context, employeeID int64) (Employee, error)
	GetJob(ctx context.Context, jobID int64) (Job, error)
	GetUser(ctx context.Context, email string) (User, error)
	ListDepartments(ctx context.Context, arg ListDepartmentsParams) ([]Department, error)
	ListEmployees(ctx context.Context, arg ListEmployeesParams) ([]Employee, error)
	ListJobs(ctx context.Context, arg ListJobsParams) ([]Job, error)
	UpdateDepartment(ctx context.Context, arg UpdateDepartmentParams) (Department, error)
	UpdateEmployee(ctx context.Context, arg UpdateEmployeeParams) (Employee, error)
	UpdateJob(ctx context.Context, arg UpdateJobParams) (Job, error)
}

var _ Querier = (*Queries)(nil)
