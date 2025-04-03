package interfaces

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"go.uber.org/zap"
)

type INamedLogger interface {
	Lgr(name string) *zap.Logger
}

type IAppContext interface {
	INamedLogger
	SM() IServiceManager
	DM() IDAOManager
	PC() IPermissionCache
}

type ISetAppContext interface {
	SetAppCtx(ctx IAppContext)
}

type IServiceContext interface {
	INamedLogger
	DM() IDAOManager
	PC() IPermissionCache
}

type IPermissionCache interface {
	AddPermission(levelID int64, perms ...model.Privileges)
	GetPermissions(levelID int64) []model.Privileges
	HasPermissionByID(levelID int64, permissionID int64) bool
	HasPermissionByName(levelID int64, permissionName string) bool
	SetPermissions(cache authModels.PermissionCache)
}
