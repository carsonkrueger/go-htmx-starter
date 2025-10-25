package model

import "github.com/carsonkrueger/main/pkg/model/db/auth"

type UserRoleJoin struct {
	auth.Users
	PLID   int64
	PLName string
}
