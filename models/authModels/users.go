package authModels

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
)

type UserPrivilegeLevelJoin struct {
	model.Users
	PLID   int64
	PLName string
}
