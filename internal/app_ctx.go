package internal

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/carsonkrueger/main/internal/services"
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

type RoutePath interface {
	Path() string
}

type PublicRoute interface {
	PublicRoute(r chi.Router)
}

type AppPublicRoute interface {
	SetAppContext
	RoutePath
	PublicRoute
}

type PrivateRoute interface {
	PrivateRoute(b *PrivateRouteBuilder)
}

type AppPrivateRoute interface {
	SetAppContext
	RoutePath
	PrivateRoute
}

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tools.AUTH_TOKEN_KEY, token)
}

func GetToken(ctx context.Context) string {
	return ctx.Value(tools.AUTH_TOKEN_KEY).(string)
}

var USER_ID_KEY = "USER_ID"

func WithUserId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, USER_ID_KEY, id)
}

func GetUserId(ctx context.Context) int64 {
	return ctx.Value(USER_ID_KEY).(int64)
}
