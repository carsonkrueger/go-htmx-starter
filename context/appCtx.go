package context

import (
	"context"
	gctx "context"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/tools"
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

func WithToken(ctx gctx.Context, token string) gctx.Context {
	return context.WithValue(ctx, tools.AUTH_TOKEN_KEY, token)
}

func GetToken(ctx gctx.Context) string {
	return ctx.Value(tools.AUTH_TOKEN_KEY).(string)
}

var USER_ID_KEY = "USER_ID"

func WithUserId(ctx gctx.Context, id int64) gctx.Context {
	return context.WithValue(ctx, USER_ID_KEY, id)
}

func GetUserId(ctx context.Context) int64 {
	return ctx.Value(USER_ID_KEY).(int64)
}
