package api

import (
	mockdb "example/employee/server/db/mock"
	db "example/employee/server/db/sqlc"
	"example/employee/server/service/role_service"
	"example/employee/server/token"
	"example/employee/server/util"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func mockAdminRole() db.Role {
	return db.Role{
		ID:         util.RandomInt64(1, 100),
		Role:       string(role_service.Admin),
		Privileges: role_service.GetAdminPrivilege(),
	}
}

func mockStaffRole() db.Role {
	return db.Role{
		ID:         util.RandomInt64(1, 100),
		Role:       string(role_service.Staff),
		Privileges: role_service.GetStaffPrivilege(),
	}
}

func mockRole(role role_service.Role) (db.Role, error) {
	fmt.Println(role)
	switch role {
	case role_service.Admin:
		return mockAdminRole(), nil
	case role_service.Staff:
		return mockStaffRole(), nil
	default:
		return mockAdminRole(), fmt.Errorf("unsupport role")
	}
}

func mockRolePrivilegeAuth(
	t *testing.T,
	store *mockdb.MockStore,
	request *http.Request,
	tokenMaker token.Maker,
	role role_service.Role,
) {
	user, _ := randomUser(t)
	mockRole, err := mockRole(role)
	require.NoError(t, err)

	token, err := tokenMaker.CreateToken(user.Email, token.Web, time.Minute)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authTypeBearer, token)
	request.Header.Set(authHeaderKey, authHeader)

	store.EXPECT().
		GetUser(gomock.Any(), gomock.Eq(user.Email)).
		Times(1).
		Return(user, nil)

	store.EXPECT().
		GetRole(gomock.Any(), gomock.Eq(user.RoleID.Int64)).
		Times(1).
		Return(mockRole, nil)
}
