package db

import (
	"context"
	"database/sql"
	"example/employee/server/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func compareTwoDate(t *testing.T, t1 time.Time, t2 time.Time) {
	// return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	require.Equal(t, t1.Year(), t2.Year())
	require.Equal(t, t1.Month(), t2.Month())
	require.Equal(t, t1.Day(), t2.Day())
}

func createRandomEmployee(t *testing.T) Employee {
	job := createRandomJob(t)
	department := createRandomDepartment(t)
	arg := CreateEmployeeParams{
		FirstName:    sql.NullString{util.RandomFirstName(), true},
		LastName:     util.RandomLastName(),
		Email:        util.RandomEmail(),
		PhoneNumber:  sql.NullString{util.RandomPhoneNumber(), true},
		HireDate:     time.Now(),
		JobID:        job.JobID,
		Salary:       util.RandomInt64(job.MinSalary.Int64, job.MaxSalary.Int64),
		ManagerID:    sql.NullInt64{0, false},
		DepartmentID: department.DepartmentID,
	}
	employee, err := testQueries.CreateEmployee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employee)

	require.Equal(t, employee.FirstName, arg.FirstName)
	require.Equal(t, employee.LastName, arg.LastName)
	require.Equal(t, employee.Email, arg.Email)
	require.Equal(t, employee.PhoneNumber, arg.PhoneNumber)
	compareTwoDate(t, employee.HireDate, arg.HireDate)
	require.Equal(t, employee.JobID, arg.JobID)
	require.Equal(t, employee.Salary, arg.Salary)
	require.Equal(t, employee.ManagerID, arg.ManagerID)
	require.Equal(t, employee.DepartmentID, arg.DepartmentID)

	require.NotZero(t, employee.EmployeeID)

	return employee
}

func TestCreateEmployee(t *testing.T) {
	createRandomEmployee(t)
}

func TestGetEmployee(t *testing.T) {
	employee1 := createRandomEmployee(t)
	employee2, err := testQueries.GetEmployee(context.Background(), employee1.EmployeeID)

	require.NoError(t, err)
	require.NotEmpty(t, employee2)

	require.Equal(t, employee1.FirstName, employee2.FirstName)
	require.Equal(t, employee1.LastName, employee2.LastName)
	require.Equal(t, employee1.Email, employee2.Email)
	require.Equal(t, employee1.PhoneNumber, employee2.PhoneNumber)
	compareTwoDate(t, employee1.HireDate, employee2.HireDate)
	require.Equal(t, employee1.JobID, employee2.JobID)
	require.Equal(t, employee1.Salary, employee2.Salary)
	require.Equal(t, employee1.ManagerID, employee2.ManagerID)
	require.Equal(t, employee1.DepartmentID, employee2.DepartmentID)
}

func TestUpdateEmployee(t *testing.T) {
	employee1 := createRandomEmployee(t)
	arg := UpdateEmployeeParams{
		EmployeeID:   employee1.EmployeeID,
		FirstName:    sql.NullString{util.RandomFirstName(), true},
		LastName:     employee1.LastName,
		Email:        util.RandomEmail(),
		PhoneNumber:  employee1.PhoneNumber,
		HireDate:     time.Now(),
		ManagerID:    employee1.ManagerID,
		DepartmentID: employee1.DepartmentID,
	}

	employee2, err := testQueries.UpdateEmployee(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, employee2)

	require.Equal(t, arg.FirstName, employee2.FirstName)
	require.Equal(t, arg.LastName, employee2.LastName)
	require.Equal(t, arg.Email, employee2.Email)
	require.Equal(t, arg.PhoneNumber, employee2.PhoneNumber)
	compareTwoDate(t, arg.HireDate, employee2.HireDate)
	require.Equal(t, arg.ManagerID, employee2.ManagerID)
	require.Equal(t, arg.DepartmentID, employee2.DepartmentID)
}

func TestDeleteEmployee(t *testing.T) {
	employee1 := createRandomEmployee(t)
	err := testQueries.DeleteEmployee(context.Background(), employee1.EmployeeID)
	require.NoError(t, err)

	employee2, err := testQueries.GetEmployee(context.Background(), employee1.EmployeeID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, employee2)

	job, err := testQueries.GetJob(context.Background(), employee1.JobID)
	require.NoError(t, err)
	require.NotEmpty(t, job)

	department, err := testQueries.GetDepartment(context.Background(), employee1.DepartmentID)
	require.NoError(t, err)
	require.NotEmpty(t, department)
}

func TestListEmployee(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEmployee(t)
	}

	arg := ListEmployeesParams{
		Limit:  5,
		Offset: 5,
	}

	employees, err := testQueries.ListEmployees(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, employees, 5)

	for _, employee := range employees {
		require.NotEmpty(t, employee)
	}
}
