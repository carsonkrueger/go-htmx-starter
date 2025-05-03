package builders

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/middlewares"
	"github.com/go-chi/chi/v5"
)

type RouteMethod string

const (
	GET    RouteMethod = "GET"
	POST   RouteMethod = "POST"
	PUT    RouteMethod = "PUT"
	PATCH  RouteMethod = "PATCH"
	DELETE RouteMethod = "DELETE"
)

type privateHandlerBuilder struct {
	appCtx         context.AppContext
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
	privDAO := mb.appCtx.DM().PrivilegeDAO()

	var r chi.Router
	if mb.permissionName != nil {
		priv := model.Privileges{Name: *mb.permissionName}
		privDAO.Upsert(&priv, mb.appCtx.DB())
		r = mb.router.With(middlewares.ApplyPermission(priv.ID, mb.appCtx))
	}
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(string(mb.method), mb.pattern, mb.handle)
}
