package builders

import (
	gctx "context"
	"net/http"
	"slices"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_starter_db/auth/model"
	"github.com/carsonkrueger/main/middlewares"
	"github.com/carsonkrueger/main/util/slice"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
	appCtx     context.AppContext
	router     chi.Router
	mw         []func(next http.Handler) http.Handler
	method     RouteMethod
	pattern    string
	handle     http.HandlerFunc
	privileges []string // privileges required to access private endpoint
}

func (mb *privateHandlerBuilder) Register(method RouteMethod, pattern string, handle http.HandlerFunc) *privateHandlerBuilder {
	mb.method = method
	mb.pattern = pattern
	mb.handle = handle
	return mb
}

// privileges required to access private endpoint
func (mb *privateHandlerBuilder) SetRequiredPrivileges(privileges []string) *privateHandlerBuilder {
	mb.privileges = privileges
	return mb
}

func (mb *privateHandlerBuilder) SetMiddlewares(middlewares ...func(next http.Handler) http.Handler) *privateHandlerBuilder {
	mb.mw = middlewares
	return mb
}

func (mb *privateHandlerBuilder) Build(ctx gctx.Context) {
	lgr := mb.appCtx.Lgr("privateHandlerBuilder.Build.")
	privDAO := mb.appCtx.DM().PrivilegeDAO()

	r := mb.router
	if len(mb.privileges) > 0 {
		privs, err := privDAO.GetManyByName(ctx, mb.privileges)
		if err != nil {
			lgr.Error("get many privileges by name", zap.Error(err))
			return
		}

		newPrivNames := slices.DeleteFunc(mb.privileges, func(privName string) bool {
			return slices.ContainsFunc(privs, func(priv model.Privileges) bool {
				return priv.Name == privName
			})
		})

		newPrivs := make([]model.Privileges, len(newPrivNames))
		for i, np := range newPrivNames {
			newPrivs[i] = model.Privileges{Name: np}
		}

		if err := privDAO.UpsertMany(ctx, newPrivs); err != nil {
			lgr.Error("upserting many privileges", zap.Error(err))
			return
		}
		privIDs := slice.MapIdx(privs, func(priv model.Privileges, _ int) int64 {
			return priv.ID
		})
		r = mb.router.With(middlewares.ApplyPermission(privIDs, mb.appCtx))
	}
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(string(mb.method), mb.pattern, mb.handle)
}
