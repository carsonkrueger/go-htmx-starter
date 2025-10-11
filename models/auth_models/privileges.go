package auth_models

import (
	"time"

	"github.com/carsonkrueger/main/gen/go_starter_db/auth/model"
)

type PermissionCache map[int64][]model.Privileges

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
	Privileges []model.Privileges
}
