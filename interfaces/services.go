package interfaces

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
)

type IServiceManager interface {
	UsersService() IUsersService
	PrivilegesService() IPrivilegesService
}

type IUsersService interface {
	Login(email string, password string) (*string, error)
	Logout(id int64, token string) error
}

type IPrivilegesService interface {
	BuildCache() error
	AddPermission(levelID int64, perms ...model.Privileges)
	GetPermissions(levelID int64) []model.Privileges
	HasPermissionByID(levelID int64, permissionID int64) bool
	HasPermissionByName(levelID int64, permissionName string) bool
}
