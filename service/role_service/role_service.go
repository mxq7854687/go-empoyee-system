package role_service

import (
	"context"
	"database/sql"
	"encoding/json"
	db "example/employee/server/db/sqlc"
	"fmt"
	"log"
)

type RoleService struct {
	Context context.Context
	Store   db.Store
}

func NewRoleService(store db.Store, context context.Context) *RoleService {
	return &RoleService{
		Context: context,
		Store:   store,
	}
}

type Role string

const (
	Admin Role = "Admin"
	Staff Role = "Staff"
)

var adminPriviledge = []db.Privilege{
	db.PrivilegeCreateAndUpdateJobs,
	db.PrivilegeCreateAndUpdateDepartments,
	db.PrivilegeDeleteJobs,
	db.PrivilegeDeleteDepartments,
	db.PrivilegeCreateAndUpdateEmployees,
	db.PrivilegeDelteEmployees,
	db.PrivilegeReadAllEmployees,
}

type RolePrivilege struct {
	pMap map[db.Privilege]bool
}

func GetAdminPrivilege() []byte {
	pMap := map[db.Privilege]bool{}
	for _, privilege := range adminPriviledge {
		pMap[privilege] = true
	}

	rolePrivilege := RolePrivilege{pMap: pMap}
	return rolePrivilege.mapToJson()
}

func GetStaffPrivilege() []byte {
	rolePrivilege := RolePrivilege{pMap: map[db.Privilege]bool{}}
	return rolePrivilege.mapToJson()
}

func (rolePrivilege *RolePrivilege) mapToJson() []byte {
	jsonByte, err := json.Marshal(rolePrivilege.pMap)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return jsonByte
}

func (roleService *RoleService) CreateAdminRole() (db.Role, error) {
	privilege := GetAdminPrivilege()
	arg := db.CreateRoleParams{
		Role:       string(Admin),
		Privileges: privilege,
	}
	return roleService.Store.CreateRole(roleService.Context, arg)
}

func (roleService *RoleService) CreateStaffRole() (db.Role, error) {
	privilege := GetStaffPrivilege()
	arg := db.CreateRoleParams{
		Role:       string(Staff),
		Privileges: privilege,
	}
	return roleService.Store.CreateRole(roleService.Context, arg)
}

func (roleService *RoleService) CreateIfNotExist(
	role Role,
	function func() (db.Role, error),
) error {
	_, err := roleService.Store.GetRoleByRoleName(roleService.Context, string(role))
	if err == sql.ErrNoRows {
		_, err = function()
	}
	return err
}

func (roleService *RoleService) InitRole() {
	err := roleService.CreateIfNotExist(Admin, roleService.CreateAdminRole)
	err = roleService.CreateIfNotExist(Staff, roleService.CreateStaffRole)
	if err != nil {
		log.Fatal("error on init role: ", err)
	}
}

func (roleService *RoleService) HasRolePriviledge(role Role, requiredPrivilege db.Privilege) error {
	getRole, err := roleService.Store.GetRoleByRoleName(roleService.Context, string(role))
	if err != nil {
		return err
	}

	var rolePrivileges map[db.Privilege]bool
	err = json.Unmarshal([]byte(getRole.Privileges), &rolePrivileges)
	if err != nil {
		return err
	}

	if _, found := rolePrivileges[requiredPrivilege]; !found {
		return fmt.Errorf("privilege %s not found", requiredPrivilege)
	}

	return nil
}
