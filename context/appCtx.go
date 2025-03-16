package context

import (
	"github.com/carsonkrueger/main/interfaces"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type appContext struct {
	lgr *zap.Logger
	sm  interfaces.IServiceManager
	dm  interfaces.IDAOManager
	pc  interfaces.IPermissionCache
}

func NewAppContext(lgr *zap.Logger, sm interfaces.IServiceManager, dm interfaces.IDAOManager, pc interfaces.IPermissionCache) *appContext {
	return &appContext{
		lgr,
		sm,
		dm,
		pc,
	}
}

func (ctx *appContext) Lgr() *zap.Logger {
	return ctx.lgr
}

func (ctx *appContext) SM() interfaces.IServiceManager {
	return ctx.sm
}

func (ctx *appContext) DM() interfaces.IDAOManager {
	return ctx.dm
}

func (ctx *appContext) PC() interfaces.IPermissionCache {
	return ctx.pc
}

func (ctx *appContext) CleanUp() {
	if err := ctx.lgr.Sync(); err != nil {
		ctx.lgr.Error("failed to sync logger", zap.Error(err))
	}
}
