package interfaces

import (
	"net/http"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
)

type IServiceManager interface {
	SetAppContext(appCtx IAppContext)
	UsersService() IUsersService
	PrivilegesService() IPrivilegesService
}

type IUsersService interface {
	Login(email string, password string, req *http.Request) (*string, error)
	Logout(id int64, token string) error
	LogoutRequest(req *http.Request) error
	GetAuthParts(req *http.Request) (string, int64, error)
}

type IPrivilegesService interface {
	BuildCache() error
	AddPermission(levelID int64, perms ...model.Privileges)
	GetPermissions(levelID int64) []model.Privileges
	HasPermissionByID(levelID int64, permissionID int64) bool
	HasPermissionByName(levelID int64, permissionName string) bool
}
