package context

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/carsonkrueger/main/templates/datadisplay"
)

type ServiceManager interface {
	// DB-START
	UsersService() UsersService
	PrivilegesService() PrivilegesService
	// DB-END
	// INSERT GET SERVICE
}

type PrivilegesService interface {
	CreatePrivilegeAssociation(ctx gctx.Context, levelID int64, privID int64) error
	DeletePrivilegeAssociation(ctx gctx.Context, levelID int64, privID int64) error
	CreateLevel(ctx gctx.Context, name string) error
	HasPermissionByID(ctx gctx.Context, levelID int64, permissionID int64) bool
	SetUserPrivilegeLevel(ctx gctx.Context, levelID int64, userID int64) error
	UserPrivilegeLevelJoinAsRowData(ctx gctx.Context, upl []auth_models.UserPrivilegeLevelJoin, allLevels []*model.PrivilegeLevels) []datadisplay.RowData
	JoinedPrivilegeLevelAsRowData(ctx gctx.Context, jpl []auth_models.JoinedPrivilegeLevel) []datadisplay.RowData
	JoinedPrivilegesAsRowData(ctx gctx.Context, jpl []auth_models.JoinedPrivilegesRaw) []datadisplay.RowData
}

type UsersService interface {
	Login(ctx gctx.Context, email string, password string, req *http.Request) (*string, error)
	Logout(ctx gctx.Context, id int64, token string) error
	LogoutRequest(ctx gctx.Context, req *http.Request) error
	GetAuthParts(ctx gctx.Context, req *http.Request) (string, int64, error)
}

// INSERT INTERFACE SERVICE
