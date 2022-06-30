// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Privilege string

const (
	PrivilegeCreateAndUpdateJobs        Privilege = "CreateAndUpdateJobs"
	PrivilegeCreateAndUpdateDepartments Privilege = "CreateAndUpdateDepartments"
	PrivilegeDeleteJobs                 Privilege = "DeleteJobs"
	PrivilegeDeleteDepartments          Privilege = "DeleteDepartments"
	PrivilegeCreateAndUpdateEmployees   Privilege = "CreateAndUpdateEmployees"
	PrivilegeDelteEmployees             Privilege = "DelteEmployees"
	PrivilegeReadAllEmployees           Privilege = "ReadAllEmployees"
)

func (e *Privilege) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Privilege(s)
	case string:
		*e = Privilege(s)
	default:
		return fmt.Errorf("unsupported scan type for Privilege: %T", src)
	}
	return nil
}

type UserStatus string

const (
	UserStatusPending     UserStatus = "pending"
	UserStatusActivated   UserStatus = "activated"
	UserStatusDeactivated UserStatus = "deactivated"
)

func (e *UserStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserStatus(s)
	case string:
		*e = UserStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UserStatus: %T", src)
	}
	return nil
}

type Department struct {
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
}

type Employee struct {
	EmployeeID   int64          `json:"employee_id"`
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

type Job struct {
	JobID     int64         `json:"job_id"`
	JobTitle  string        `json:"job_title"`
	MinSalary sql.NullInt64 `json:"min_salary"`
	MaxSalary sql.NullInt64 `json:"max_salary"`
}

type Role struct {
	ID         int64           `json:"id"`
	Role       string          `json:"role"`
	Privileges json.RawMessage `json:"privileges"`
	UpdatedAt  time.Time       `json:"updated_at"`
	CreatedAt  time.Time       `json:"created_at"`
}

type User struct {
	Email          string     `json:"email"`
	Status         UserStatus `json:"status"`
	HashedPassword string     `json:"hashed_password"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedAt      time.Time  `json:"created_at"`
}
