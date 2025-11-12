package middlewares

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/carsonkrueger/main/pkg/util"
	"go.uber.org/zap"
)

func Auth(appCtx *context.AppContext, enforce bool) func(next http.Handler) http.Handler {
	usersService := appCtx.SM().UsersService()
	usersDAO := appCtx.DM().UsersDAO()
	sessionsDAO := appCtx.DM().SessionsDAO()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			lgr := context.GetLogger(ctx, "mw.Auth")

			token, id, err := usersService.GetAuthParts(ctx, req)
			if err != nil {
				if enforce {
					util.HandleRedirect(req, w, "/login")
					common.HandleError(req, w, lgr, err, 403, "Malformed auth token")
					return
				}
			}

			user, err := usersDAO.GetOne(ctx, id)
			if err != nil {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				if enforce {
					util.HandleRedirect(req, w, "/login")
					common.HandleError(req, w, lgr, err, 403, "Malformed auth token")
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
					util.HandleRedirect(req, w, "/login")
					common.HandleError(req, w, lgr, err, 403, "Malformed auth token")
					return
				}
			}

			if session.Token != token {
				req.Header.Del(constant.AUTH_TOKEN_KEY)
				if enforce {
					util.HandleRedirect(req, w, "/login")
					common.HandleError(req, w, lgr, err, 403, "Malformed auth token")
					return
				}
			}

			ctx = context.WithToken(ctx, token)
			ctx = context.WithUserId(ctx, id)
			ctx = context.WithRoleID(ctx, user.RoleID)
			ctx = context.WithLogger(ctx, lgr.With(zap.Int64("user_id", id)))

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
