package middlewares

import (
	"net/http"

	"github.com/carsonkrueger/main/constant"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/carsonkrueger/main/tools"
)

func EnforceAuth(appCtx context.AppContext) func(next http.Handler) http.Handler {
	usersService := appCtx.SM().UsersService()
	usersDAO := appCtx.DM().UsersDAO()
	sessionsDAO := appCtx.DM().SessionsDAO()
	lgr := appCtx.Lgr("MW EnforceAuth")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			token, id, err := usersService.GetAuthParts(ctx, req)
			if err != nil {
				tools.HandleError(req, res, lgr, err, 403, "Malformed auth token")
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			user, err := usersDAO.GetOne(ctx, id)
			if err != nil {
				tools.HandleError(req, res, lgr, err, 403, "Malformed auth token")
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			key := auth_models.SessionsPrimaryKey{
				UserID:    id,
				AuthToken: token,
			}
			session, err := sessionsDAO.GetOne(ctx, key)

			if err != nil {
				tools.HandleError(req, res, lgr, err, 403, "Malformed auth token")
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			if session.Token != token {
				tools.HandleError(req, res, lgr, err, 403, "Malformed auth token")
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			ctx = context.WithToken(ctx, token)
			ctx = context.WithUserId(ctx, id)
			ctx = context.WithPrivilegeLevelID(ctx, user.PrivilegeLevelID)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
