package builders

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/go-chi/chi/v5"
)

type IRoutePath interface {
	Path() string
}

type IPublicRoute interface {
	PublicRoute(r chi.Router)
}

type IAppPublicRoute interface {
	IRoutePath
	IPublicRoute
}

// DB-START
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

type IPrivateRoute interface {
	PrivateRoute(b *PrivateRouteBuilder)
}

type IAppPrivateRoute interface {
	IRoutePath
	IPrivateRoute
}

func (rb *PrivateRouteBuilder) NewHandle() *privateHandlerBuilder {
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

func (rb *PrivateRouteBuilder) AddMiddleware(middleware func(next http.Handler) http.Handler) {
	rb.router.Use(middleware)
}

func (rb *PrivateRouteBuilder) RawRouter() chi.Router {
	return rb.router
}

// DB-END
