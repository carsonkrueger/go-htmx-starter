package authModels

import (
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
)

type PermissionCache map[int64][]model.Privileges
type LevelNameCache map[string][]int64

type PrivilegeLevelsPrivilegesPrimaryKey struct {
	PrivilegeID      int64
	PrivilegeLevelID int64
}

type JoinedPrivilegesRaw struct {
	LevelID            int64
	LevelName          string
	PrivilegeID        int64
	PrivilegeName      string
	PrivilegeCreatedAt *time.Time
}

type JoinedPrivilegeLevel struct {
	LevelID    int64
	LevelName  string
	Privileges []model.Privileges
}
