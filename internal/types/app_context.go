package types

import (
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext struct {
	Lgr *zap.Logger
	Sm  IServiceManager
}

func NewAppContext(lgr *zap.Logger, sm IServiceManager) *AppContext {
	return &AppContext{
		lgr,
		sm,
	}
}

func (ctx *AppContext) CleanUp() {
	if err := ctx.Lgr.Sync(); err != nil {
		ctx.Lgr.Error("failed to sync logger", zap.Error(err))
	}
}

type SetAppContext interface {
	SetAppCtx(ctx *AppContext)
}

type WithAppContext struct {
	AppCtx *AppContext
}

func (b *WithAppContext) SetAppCtx(ctx *AppContext) {
	b.AppCtx = ctx
}
