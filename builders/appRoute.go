package builders

import (
	"net/http"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/middlewares"
	"github.com/go-chi/chi/v5"
)

type RouteMethod string

const (
	GET    RouteMethod = "GET"
	POST   RouteMethod = "POST"
	PUT    RouteMethod = "PUT"
	DELETE RouteMethod = "DELETE"
)

type RoutePath interface {
	Path() string
}

type PublicRoute interface {
	PublicRoute(r chi.Router)
}

type AppPublicRoute interface {
	interfaces.ISetAppContext
	RoutePath
	PublicRoute
}

type PrivateRoute interface {
	PrivateRoute(b *PrivateRouteBuilder)
}

type AppPrivateRoute interface {
	interfaces.ISetAppContext
	RoutePath
	PrivateRoute
}

type PrivateRouteBuilder struct {
	router chi.Router
	appCtx interfaces.IAppContext
}

type privateMethodBuilder struct {
	appCtx         interfaces.IAppContext
	router         chi.Router
	mw             []func(next http.Handler) http.Handler
	method         RouteMethod
	pattern        string
	handle         http.HandlerFunc
	permissionName *string
}

func NewPrivateRouteBuilder(appCtx interfaces.IAppContext) PrivateRouteBuilder {
	return PrivateRouteBuilder{
		router: chi.NewRouter(),
		appCtx: appCtx,
	}
}

func (rb *PrivateRouteBuilder) NewHandle() *privateMethodBuilder {
	return &privateMethodBuilder{
		router: rb.router,
		appCtx: rb.appCtx,
	}
}

func (mb *PrivateRouteBuilder) Build() chi.Router {
	return mb.router
}

func (mb *privateMethodBuilder) RegisterRoute(method RouteMethod, pattern string, handle http.HandlerFunc) *privateMethodBuilder {
	mb.method = method
	mb.pattern = pattern
	mb.handle = handle
	return mb
}

func (mb *privateMethodBuilder) SetPermissionName(name string) *privateMethodBuilder {
	mb.permissionName = &name
	return mb
}

func (mb *privateMethodBuilder) SetMiddlewares(middlewares ...func(next http.Handler) http.Handler) *privateMethodBuilder {
	mb.mw = middlewares
	return mb
}

func (mb *privateMethodBuilder) Build() {
	var r chi.Router
	if mb.permissionName != nil {
		r = mb.router.With(middlewares.ApplyPermission(*mb.permissionName, mb.appCtx))
	}
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(string(mb.method), mb.pattern, mb.handle)
}
