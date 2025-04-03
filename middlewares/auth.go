package middlewares

import (
	"errors"
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/tools"
)

func EnforceAuth(appCtx interfaces.IAppContext) func(next http.Handler) http.Handler {
	dao := appCtx.DM().UsersDAO()
	lgr := appCtx.Lgr("MW EnforceAuth")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			cookie, err := tools.GetAuthCookie(req)
			if err != nil {
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Not Authenticated"))
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			token, id, err := tools.GetAuthParts(cookie)
			if err != nil {
				req.Header.Del(tools.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, err)
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			user, err := dao.GetById(id)

			if err != nil {
				req.Header.Del(tools.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Malformed auth token"))
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			if *user.AuthToken != cookie.Value {
				req.Header.Del(tools.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Malformed auth token"))
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			ctx = context.WithToken(ctx, token)
			ctx = context.WithUserId(ctx, id)
			ctx = context.WithPrivilegeID(ctx, user.PrivilegeLevelID)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
