package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/carsonkrueger/main/internal/context"
	"go.uber.org/zap"
)

func Recover(appCtx context.AppContext) func(next http.Handler) http.Handler {
	lgr := appCtx.Lgr("NoAuth")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					lgr.Error("panic", zap.String("stack", string(stack)))
					res.WriteHeader(http.StatusInternalServerError)
				}
			}()
			ctx := req.Context()
			ctx = context.WithDB(ctx, appCtx.DB())

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
