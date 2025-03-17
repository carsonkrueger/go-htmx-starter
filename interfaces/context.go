package interfaces

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"go.uber.org/zap"
)

type IAppContext interface {
	Lgr() *zap.Logger
	SM() IServiceManager
	DM() IDAOManager
	PC() IPermissionCache
}

type ISetAppContext interface {
	SetAppCtx(ctx IAppContext)
}

type IServiceContext interface {
	Lgr() *zap.Logger
	DM() IDAOManager
	PC() IPermissionCache
}

type IPermissionCache interface {
	AddPermission(levelID int64, perm model.Privileges)
	SetPermissions(map[int64][]model.Privileges)
	GetPermissions(levelID int64) []model.Privileges
	HasPermissionByID(levelID int64, permissionID int64) bool
	HasPermissionByName(levelID int64, permissionName string) bool
}
