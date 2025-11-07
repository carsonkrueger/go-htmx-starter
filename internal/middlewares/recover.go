package middlewares

import (
	gctx "context"
	"net/http"
	"runtime/debug"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Recover(ctx gctx.Context, appCtx *context.AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			reqCtx := req.Context()
			reqID, _ := uuid.NewV7()
			lgr := context.GetLogger(ctx).With(zap.String("req_id", reqID.String()))

			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					lgr.Error("panic", zap.Any("recover_err:", err), zap.String("stack", string(stack)))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			reqCtx = context.WithLogger(reqCtx, lgr)
			reqCtx = context.WithDB(reqCtx, context.GetDB(ctx))

			next.ServeHTTP(w, req.WithContext(reqCtx))
		})
	}
}
