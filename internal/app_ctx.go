package internal

import (
	"github.com/go-chi/chi/v5"

	"github.com/carsonkrueger/main/internal/services"
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
