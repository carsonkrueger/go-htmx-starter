package middlewares

import (
	"errors"
	"net/http"

	"github.com/carsonkrueger/main/constant"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/tools"
)

func EnforceAuth(appCtx interfaces.IAppContext) func(next http.Handler) http.Handler {
	usersService := appCtx.SM().UsersService()
	usersDAO := appCtx.DM().UsersDAO()
	sessionsDAO := appCtx.DM().SessionsDAO()
	lgr := appCtx.Lgr("MW EnforceAuth")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			token, id, err := usersService.GetAuthParts(req)
			if err != nil {
				tools.RequestHttpError(ctx, lgr, res, 403, err)
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			user, err := usersDAO.GetById(id)
			if err != nil {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Malformed auth token"))
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			key := authModels.SessionsPrimaryKey{
				UserID:    id,
				AuthToken: token,
			}
			session, err := sessionsDAO.GetById(key)

			if session == nil || err != nil {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Malformed auth token"))
				res.Header().Set("Hx-Redirect", "/login")
				return
			}

			if session.Token != token {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				tools.RequestHttpError(ctx, lgr, res, 403, errors.New("Malformed auth token"))
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
