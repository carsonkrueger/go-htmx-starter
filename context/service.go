package context

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/gen/go_starter_db/auth/model"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/carsonkrueger/main/templates/datadisplay"
)

type ServiceManager interface {
	UsersService() UsersService
	PrivilegesService() PrivilegesService
	// INSERT GET SERVICE
}

type PrivilegesService interface {
	CreatePrivilegeAssociation(ctx gctx.Context, role int16, privID int64) error
	DeletePrivilegeAssociation(ctx gctx.Context, role int16, privID int64) error
	CreateRole(ctx gctx.Context, name string) error
	HasPermissionsByIDS(ctx gctx.Context, role int16, privileges []int64) bool
	SetUserRole(ctx gctx.Context, role int16, userID int64) error
	UserRoleJoinAsRowData(ctx gctx.Context, upl []auth_models.UserRoleJoin, roles []model.Roles) []datadisplay.RowData
	JoinedRoleAsRowData(ctx gctx.Context, jpl []auth_models.JoinedRole) []datadisplay.RowData
	JoinedPrivilegesAsRowData(ctx gctx.Context, jpl []auth_models.JoinedPrivilegesRaw) []datadisplay.RowData
}

type UsersService interface {
	Login(ctx gctx.Context, email string, password string, req *http.Request) (*string, error)
	Logout(ctx gctx.Context, id int64, token string) error
	LogoutRequest(ctx gctx.Context, req *http.Request) error
	GetAuthParts(ctx gctx.Context, req *http.Request) (string, int64, error)
}

// INSERT INTERFACE SERVICE
