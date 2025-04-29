package builders

import (
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

func (mb *PrivateRouteBuilder) Build() chi.Router {
	return mb.router
}

// DB-END
