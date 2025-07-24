package middlewares

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
)

func NoAuth(appCtx context.AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			db := appCtx.DB()
			ctx = context.WithDB(ctx, db)
			req.WithContext(ctx)
			next.ServeHTTP(res, req)
		})
	}
}
