package builders

import (
	gctx "context"
	"net/http"
	"slices"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/middlewares"
	"github.com/carsonkrueger/main/pkg/model/db/auth"
	"github.com/carsonkrueger/main/pkg/util/slice"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type privateHandlerBuilder struct {
	appCtx     context.AppContext
	router     chi.Router
	mw         []func(next http.Handler) http.Handler
	method     string
	pattern    string
	handle     http.HandlerFunc
	privileges []string // privileges required to access private endpoint
}

func (mb *privateHandlerBuilder) Register(method string, pattern string, handle http.HandlerFunc) *privateHandlerBuilder {
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
	lgr := mb.appCtx.Lgr("privateHandlerBuilder.Build")
	privDAO := mb.appCtx.DM().PrivilegeDAO()

	r := mb.router
	if len(mb.privileges) > 0 {
		privs, err := privDAO.GetManyByName(ctx, mb.privileges)
		if err != nil {
			lgr.Error("get many privileges by name", zap.Error(err))
			return
		}

		newPrivNames := slices.DeleteFunc(mb.privileges, func(privName string) bool {
			return slices.ContainsFunc(privs, func(priv auth.Privileges) bool {
				return priv.Name == privName
			})
		})

		newPrivs := make([]auth.Privileges, len(newPrivNames))
		for i, np := range newPrivNames {
			newPrivs[i] = auth.Privileges{Name: np}
		}

		if err := privDAO.UpsertMany(ctx, newPrivs); err != nil {
			lgr.Error("upserting many privileges", zap.Error(err))
			return
		}
		privIDs := slice.MapIdx(privs, func(priv auth.Privileges, _ int) int64 {
			return priv.ID
		})
		r = mb.router.With(middlewares.ApplyPermission(privIDs, mb.appCtx))
	}
	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}
	r.MethodFunc(string(mb.method), mb.pattern, mb.handle)
}
