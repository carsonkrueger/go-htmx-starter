package model

import "github.com/carsonkrueger/main/pkg/db/auth/model"

type UserRoleJoin struct {
	model.Users
	PLID   int64
	PLName string
}
