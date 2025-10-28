package model

import (
	"github.com/carsonkrueger/main/pkg/db/auth/model"
)

type PermissionCache map[int64][]model.Privileges

type RolesPrivilegesPrimaryKey struct {
	PrivilegeID int64
	RoleID      int16
}

type RolesPrivilegeJoin struct {
	model.Privileges
	model.Roles
	// RoleID             int16
	// RoleName           string
	// PrivilegeID        int64
	// PrivilegeName      string
	// PrivilegeCreatedAt *time.Time
}

type JoinedRole struct {
	RoleID     int16
	RoleName   string
	Privileges []model.Privileges
}
