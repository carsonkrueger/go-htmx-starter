package model

import "github.com/carsonkrueger/main/pkg/db/auth/model"

type UserRoleJoin struct {
	model.Users
	model.Roles
	// RoleID   int64
	// RoleName string
}
