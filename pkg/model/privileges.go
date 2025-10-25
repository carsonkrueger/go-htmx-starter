package model

import (
	"time"

	"github.com/carsonkrueger/main/pkg/model/db/auth"
)

type PermissionCache map[int64][]auth.Privileges

type RolesPrivilegesPrimaryKey struct {
	PrivilegeID int64
	RoleID      int16
}

type JoinedPrivilegesRaw struct {
	RoleID             int16
	RoleName           string
	PrivilegeID        int64
	PrivilegeName      string
	PrivilegeCreatedAt *time.Time
}

type JoinedRole struct {
	RoleID     int16
	RoleName   string
	Privileges []auth.Privileges
}
