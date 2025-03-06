package internal

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/enums"
	"github.com/go-chi/chi/v5"
)

type PrivateRouteBuilder struct {
	router chi.Router
}

func NewPrivateRouteBuilder() PrivateRouteBuilder {
	return PrivateRouteBuilder{
		router: chi.NewRouter(),
	}
}

func (rb *PrivateRouteBuilder) NewHandle() *privateMethodBuilder {
	return &privateMethodBuilder{
		router: rb.router,
	}
}

func (mb *PrivateRouteBuilder) Build() chi.Router {
	return mb.router
}

type privateMethodBuilder struct {
	router     chi.Router
	mw         []func(next http.Handler) http.Handler
	method     string
	pattern    string
	handle     http.HandlerFunc
	permission enums.Permission
}

func (mb *privateMethodBuilder) RegisterRoute(method string, pattern string, handle http.HandlerFunc) *privateMethodBuilder {
	mb.method = method
	mb.pattern = pattern
	mb.handle = handle
	return mb
}

func (mb *privateMethodBuilder) SetPermission(permission enums.Permission) *privateMethodBuilder {
	mb.permission = permission
	return mb
}

func (mb *privateMethodBuilder) SetMiddlewares(middlewares ...func(next http.Handler) http.Handler) *privateMethodBuilder {
	mb.mw = middlewares
	return mb
}

func (mb *privateMethodBuilder) Build() {
	var r chi.Router
	if mb.permission != "" {
		r = mb.router.With(ApplyPermission(mb.permission))
	}
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(mb.method, mb.pattern, mb.handle)
}
