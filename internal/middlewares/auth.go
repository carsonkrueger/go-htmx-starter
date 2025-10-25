package middlewares

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/carsonkrueger/main/pkg/util"
)

func Auth(appCtx context.AppContext, enforce bool) func(next http.Handler) http.Handler {
	usersService := appCtx.SM().UsersService()
	usersDAO := appCtx.DM().UsersDAO()
	sessionsDAO := appCtx.DM().SessionsDAO()
	lgr := appCtx.Lgr("MW EnforceAuth")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			token, id, err := usersService.GetAuthParts(ctx, req)
			if err != nil {
				if enforce {
					util.HandleError(req, res, lgr, err, 403, "Malformed auth token")
					res.Header().Set("Hx-Redirect", "/login")
					return
				}
			}

			user, err := usersDAO.GetOne(ctx, id)
			if err != nil {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				if enforce {
					util.HandleError(req, res, lgr, err, 403, "Malformed auth token")
					res.Header().Set("Hx-Redirect", "/login")
					return
				}
			}

			key := model.SessionsPrimaryKey{
				UserID:    id,
				AuthToken: token,
			}
			session, err := sessionsDAO.GetOne(ctx, key)

			if err != nil {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				if enforce {
					util.HandleError(req, res, lgr, err, 403, "Malformed auth token")
					res.Header().Set("Hx-Redirect", "/login")
					return
				}
			}

			if session.Token != token {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				if enforce {
					util.HandleError(req, res, lgr, err, 403, "Malformed auth token")
					res.Header().Set("Hx-Redirect", "/login")
					return
				}
			}

			ctx = context.WithToken(ctx, token)
			ctx = context.WithUserId(ctx, id)
			ctx = context.WithRoleID(ctx, user.RoleID)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
