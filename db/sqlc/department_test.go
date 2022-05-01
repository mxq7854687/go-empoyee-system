package db

import (
	"context"
	"database/sql"
	"example/employee/server/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomDepartment(t *testing.T) Department {
	deparmentName := util.RandomDepartmentName()

	department, err := testQueries.CreateDepartment(context.Background(), deparmentName)
	require.NoError(t, err)
	require.NotEmpty(t, department)

	require.Equal(t, deparmentName, department.DepartmentName)

	require.NotZero(t, department.DepartmentID)

	return department
}

func TestCreateDepartment(t *testing.T) {
	createRandomDepartment(t)
}

func TestGetDepartment(t *testing.T) {
	department1 := createRandomDepartment(t)
	department2, err := testQueries.GetDepartment(context.Background(), department1.DepartmentID)

	require.NoError(t, err)
	require.NotEmpty(t, department2)

	require.Equal(t, department1.DepartmentID, department2.DepartmentID)
	require.Equal(t, department1.DepartmentName, department2.DepartmentName)
}

func TestUpdateDepartment(t *testing.T) {
	department1 := createRandomDepartment(t)
	arg := UpdateDepartmentParams{
		DepartmentID:   department1.DepartmentID,
		DepartmentName: util.RandomDepartmentName(),
	}

	department2, err := testQueries.UpdateDepartment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, department2)

	require.Equal(t, department1.DepartmentID, department2.DepartmentID)
	require.Equal(t, arg.DepartmentName, department2.DepartmentName)
}

func TestDeleteDepartment(t *testing.T) {
	department1 := createRandomDepartment(t)
	err := testQueries.DeleteDepartments(context.Background(), department1.DepartmentID)
	require.NoError(t, err)

	department2, err := testQueries.GetDepartment(context.Background(), department1.DepartmentID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, department2)
}

func TestListDepartment(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomDepartment(t)
	}

	arg := ListDepartmentsParams{
		Limit:  5,
		Offset: 5,
	}

	departments, err := testQueries.ListDepartments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, departments, 5)

	for _, department := range departments {
		require.NotEmpty(t, department)
	}
}
