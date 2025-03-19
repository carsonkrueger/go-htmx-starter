package authModels

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models"
)

type PermissionCache []models.Pair[int64, []model.Privileges]

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
