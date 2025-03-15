package middlewares

import (
	"errors"
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/tools"
)

func EnforceAuth(appCtx context.IAppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			lgr := appCtx.Lgr()
			ctx := req.Context()
			cookie, err := tools.GetAuthCookie(req)
			if err != nil {
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Not Authenticated"))
				return
			}

			token, id, err := tools.GetAuthParts(cookie)
			if err != nil {
				req.Header.Del(tools.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, err)
				return
			}

			context.WithToken(ctx, token)
			context.WithUserId(ctx, id)

			dao := appCtx.DM().UsersDAO()
			user, err := dao.GetById(id)

			if err != nil {
				req.Header.Del(tools.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Malformed auth token"))
				return
			}
			if *user.AuthToken != cookie.Value {
				req.Header.Del(tools.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Malformed auth token"))
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}
