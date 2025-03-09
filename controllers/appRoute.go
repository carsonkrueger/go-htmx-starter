package controllers

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/middlewares"
	"github.com/go-chi/chi/v5"
)

type RoutePath interface {
	Path() string
}

type PublicRoute interface {
	PublicRoute(r chi.Router)
}

type AppPublicRoute interface {
	context.SetAppContext
	RoutePath
	PublicRoute
}

type PrivateRoute interface {
	PrivateRoute(b *PrivateRouteBuilder)
}

type AppPrivateRoute interface {
	context.SetAppContext
	RoutePath
	PrivateRoute
}

type PrivateRouteBuilder struct {
	router chi.Router
	appCtx context.IAppContext
}

func NewPrivateRouteBuilder(appCtx context.IAppContext) PrivateRouteBuilder {
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

type privateMethodBuilder struct {
	appCtx     context.IAppContext
	router     chi.Router
	mw         []func(next http.Handler) http.Handler
	method     string
	pattern    string
	handle     http.HandlerFunc
	permission *model.Privileges
}

func (mb *privateMethodBuilder) RegisterRoute(method string, pattern string, handle http.HandlerFunc) *privateMethodBuilder {
	mb.method = method
	mb.pattern = pattern
	mb.handle = handle
	return mb
}

func (mb *privateMethodBuilder) SetPermission(permission *model.Privileges) *privateMethodBuilder {
	mb.permission = permission
	return mb
}

func (mb *privateMethodBuilder) SetMiddlewares(middlewares ...func(next http.Handler) http.Handler) *privateMethodBuilder {
	mb.mw = middlewares
	return mb
}

func (mb *privateMethodBuilder) Build() {
	var r chi.Router
	if mb.permission != nil {
		r = mb.router.With(middlewares.ApplyPermission(mb.permission, mb.appCtx))
	}
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(mb.method, mb.pattern, mb.handle)
}
