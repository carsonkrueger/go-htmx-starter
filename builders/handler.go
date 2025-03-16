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

type privateHandlerBuilder struct {
	appCtx         interfaces.IAppContext
	router         chi.Router
	mw             []func(next http.Handler) http.Handler
	method         RouteMethod
	pattern        string
	handle         http.HandlerFunc
	permissionName *string
}

func (mb *privateHandlerBuilder) Register(method RouteMethod, pattern string, handle http.HandlerFunc) *privateHandlerBuilder {
	mb.method = method
	mb.pattern = pattern
	mb.handle = handle
	return mb
}

func (mb *privateHandlerBuilder) SetPermissionName(name string) *privateHandlerBuilder {
	mb.permissionName = &name
	return mb
}

func (mb *privateHandlerBuilder) SetMiddlewares(middlewares ...func(next http.Handler) http.Handler) *privateHandlerBuilder {
	mb.mw = middlewares
	return mb
}

func (mb *privateHandlerBuilder) Build() {
	var r chi.Router
	if mb.permissionName != nil {
		r = mb.router.With(middlewares.ApplyPermission(*mb.permissionName, mb.appCtx))
	}
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(string(mb.method), mb.pattern, mb.handle)
}
