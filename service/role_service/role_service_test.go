package role_service

import (
	"context"
	db "example/employee/server/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAdminRole(t *testing.T) {
	role, err := roleService.CreateAdminRole()
	require.NoError(t, err)
	require.Equal(t, string(Admin), role.Role)
	require.JSONEq(t, string(GetAdminPrivilege()), string(role.Privileges))
}

func TestCreateStaffRole(t *testing.T) {
	role, err := roleService.CreateStaffRole()
	require.NoError(t, err)
	require.Equal(t, string(Staff), role.Role)
	require.JSONEq(t, string(GetStaffPrivilege()), string(role.Privileges))
}

func TestInitRole(t *testing.T) {
	roleService.InitRole()

	adminRole, err := roleService.Store.GetRoleByRoleName(context.Background(), string(Admin))
	require.NoError(t, err)
	require.Equal(t, string(Admin), adminRole.Role)

	staffRole, err := roleService.Store.GetRoleByRoleName(context.Background(), string(Staff))
	require.NoError(t, err)
	require.Equal(t, string(Staff), staffRole.Role)
}

func TestHasRolePriviledge(t *testing.T) {
	_, err := roleService.CreateStaffRole()
	require.NoError(t, err)
	_, err = roleService.CreateAdminRole()
	require.NoError(t, err)

	err = roleService.HasRolePriviledge(Admin, db.PrivilegeCreateAndUpdateDepartments)
	require.NoError(t, err)

	err = roleService.HasRolePriviledge(Staff, db.PrivilegeCreateAndUpdateDepartments)
	require.Error(t, err)
}
