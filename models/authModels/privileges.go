package authModels

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
)

type PermissionCache map[int64][]model.Privileges

type JoinedPrivilegesRaw struct {
	LevelID    int64
	LevelName  string
	Privileges model.Privileges
}

type JoinedPrivilegeLevel struct {
	LevelID    int64
	LevelName  string
	Privileges []model.Privileges
}
