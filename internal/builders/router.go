package builders

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/go-chi/chi/v5"
)

type RoutePath interface {
	Path() string
}

type PublicRoute interface {
	PublicRoute(r chi.Router)
}

type AppPublicRoute interface {
	RoutePath
	PublicRoute
}

type PrivateRouteBuilder struct {
	router chi.Router
	appCtx context.AppContext
}

func NewPrivateRouteBuilder(appCtx context.AppContext) PrivateRouteBuilder {
	return PrivateRouteBuilder{
		router: chi.NewRouter(),
		appCtx: appCtx,
	}
}

type PrivateRoute interface {
	PrivateRoute(ctx gctx.Context, b *PrivateRouteBuilder)
}

type AppPrivateRoute interface {
	RoutePath
	PrivateRoute
}

func (rb *PrivateRouteBuilder) NewHandler() *privateHandlerBuilder {
	return &privateHandlerBuilder{
		router: rb.router,
		appCtx: rb.appCtx,
	}
}

func (rb *PrivateRouteBuilder) NewGroup(f func(g *PrivateRouteBuilder)) {
	builder := PrivateRouteBuilder{
		router: nil,
		appCtx: rb.appCtx,
	}
	rb.router.Group(func(g chi.Router) {
		builder.router = g
	})
	f(&builder)
}

func (rb *PrivateRouteBuilder) NewRoute(path string) PrivateRouteBuilder {
	builder := PrivateRouteBuilder{
		router: nil,
		appCtx: rb.appCtx,
	}
	rb.router.Route(path, func(r chi.Router) {
		builder.router = r
	})
	return builder
}

func (rb *PrivateRouteBuilder) AddMiddleware(middleware func(next http.Handler) http.Handler) {
	rb.router.Use(middleware)
}

func (rb *PrivateRouteBuilder) RawRouter() chi.Router {
	return rb.router
}
