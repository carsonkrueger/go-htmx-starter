package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/carsonkrueger/main/context"
	"go.uber.org/zap"
)

func NoAuth(appCtx context.AppContext) func(next http.Handler) http.Handler {
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
			// DB-START
			ctx := req.Context()
			db := appCtx.DB()
			ctx = context.WithDB(ctx, db)
			req = req.WithContext(ctx)
			// DB-END
			next.ServeHTTP(res, req)
		})
	}
}
