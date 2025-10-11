package auth_models

import (
	"github.com/carsonkrueger/main/gen/go_starter_db/auth/model"
)

type UserRoleJoin struct {
	model.Users
	PLID   int64
	PLName string
}
