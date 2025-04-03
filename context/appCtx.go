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
}

func NewAppContext(lgr *zap.Logger, sm interfaces.IServiceManager, dm interfaces.IDAOManager) *appContext {
	return &appContext{
		lgr,
		sm,
		dm,
	}
}

func (ctx *appContext) Lgr(name string) *zap.Logger {
	return ctx.lgr.Named(name)
}

func (ctx *appContext) SM() interfaces.IServiceManager {
	return ctx.sm
}

func (ctx *appContext) DM() interfaces.IDAOManager {
	return ctx.dm
}

func (ctx *appContext) CleanUp() {
	if err := ctx.lgr.Sync(); err != nil {
		ctx.lgr.Error("failed to sync logger", zap.Error(err))
	}
}
