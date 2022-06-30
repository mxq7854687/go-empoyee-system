package db

import (
	"context"
	"encoding/json"
	"example/employee/server/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomRole(t *testing.T) Role {
	ans, err := json.Marshal("{}")
	require.NoError(t, err)

	arg := CreateRoleParams{
		Role:       util.RandomString(6),
		Privileges: ans,
	}
	role, err := testQueries.CreateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role)

	require.Equal(t, role.Role, arg.Role)
	require.Equal(t, role.Privileges, arg.Privileges)

	require.NotZero(t, role.CreatedAt)
	require.NotZero(t, role.UpdatedAt)

	return role
}

func TestCreateRole(t *testing.T) {
	createRandomRole(t)
}

func TestGetRole(t *testing.T) {
	role1 := createRandomRole(t)
	role2, err := testQueries.GetRole(context.Background(), role1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.Role, role2.Role)
	require.Equal(t, role1.Privileges, role2.Privileges)

	require.WithinDuration(t, role1.UpdatedAt, role2.UpdatedAt, time.Second)
	require.WithinDuration(t, role1.CreatedAt, role2.CreatedAt, time.Second)
}

func TestGetRoleByName(t *testing.T) {
	role := createRandomRole(t)

	getRole, err := testQueries.GetRoleByRoleName(context.Background(), role.Role)
	require.NoError(t, err)
	require.Equal(t, role.ID, getRole.ID)
	require.Equal(t, role.Role, getRole.Role)
	require.Equal(t, role.Privileges, getRole.Privileges)
}
