package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/tools"
)

func ApplyPermission(p *model.Privileges, appCtx context.IAppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			lgr := appCtx.Lgr()
			lgr.Info(fmt.Sprintf("Permission Auth: %+v", p))
			cookie, err := tools.GetAuthCookie(req)
			if err != nil {
				tools.RequestHttpError(appCtx.Lgr(), res, 403, errors.New("invalid cookie"))
				return
			}

			_, id, err := tools.GetAuthParts(cookie)
			if err != nil {
				tools.RequestHttpError(appCtx.Lgr(), res, 403, err)
				return
			}
			lgr.Info(fmt.Sprintf("User id: %d", id))

			us := appCtx.SM().UsersService()
			permitted := us.IsPermitted(id, p.ID)
			if !permitted {
				tools.RequestHttpError(appCtx.Lgr(), res, 403, errors.New("permission denied"))
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
