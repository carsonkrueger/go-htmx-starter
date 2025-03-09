package context

import (
	"context"
	gctx "context"

	"github.com/carsonkrueger/main/services"
	"github.com/carsonkrueger/main/tools"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext struct {
	Lgr *zap.Logger
	SM  *services.ServiceManager
}

func NewAppContext(lgr *zap.Logger, sm *services.ServiceManager) *AppContext {
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
