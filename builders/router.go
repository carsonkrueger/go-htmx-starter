package builders

import (
	"github.com/carsonkrueger/main/interfaces"
	"github.com/go-chi/chi/v5"
)

type PrivateRouteBuilder struct {
	router chi.Router
	appCtx interfaces.IAppContext
}

func NewPrivateRouteBuilder(appCtx interfaces.IAppContext) PrivateRouteBuilder {
	return PrivateRouteBuilder{
		router: chi.NewRouter(),
		appCtx: appCtx,
	}
}

type IRoutePath interface {
	Path() string
}

type IPublicRoute interface {
	PublicRoute(r chi.Router)
}

type IAppPublicRoute interface {
	interfaces.ISetAppContext
	IRoutePath
	IPublicRoute
}

type IPrivateRoute interface {
	PrivateRoute(b *PrivateRouteBuilder)
}

type IAppPrivateRoute interface {
	interfaces.ISetAppContext
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
