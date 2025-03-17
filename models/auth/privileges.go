package auth

import "github.com/carsonkrueger/main/gen/go_db/auth/model"

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
