package builders

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/enums"
	"github.com/carsonkrueger/main/internal/middlewares"
	// "github.com/carsonkrueger/main/internal/types"
	"github.com/go-chi/chi/v5"
)

type privateMethodBuilder struct {
	router     chi.Router
	mw         []func(next http.Handler) http.Handler
	method     string
	pattern    string
	handle     http.HandlerFunc
	permission enums.Permission
}

func (mb *privateMethodBuilder) RegisterRoute(method string, pattern string, handle http.HandlerFunc, permission enums.Permission) *privateMethodBuilder {
	mb.method = method
	mb.pattern = pattern
	mb.handle = handle
	mb.permission = permission
	return mb
}

func (mb *privateMethodBuilder) SetMiddlewares(middlewares ...func(next http.Handler) http.Handler) *privateMethodBuilder {
	mb.mw = middlewares
	return mb
}

func (mb *privateMethodBuilder) Build() {
	r := mb.router.With(middlewares.ApplyPermission(mb.permission))
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(mb.method, mb.pattern, mb.handle)
}
