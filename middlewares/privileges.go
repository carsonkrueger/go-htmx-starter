package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/tools"
)

func ApplyPermission(permissionName string, appCtx interfaces.IAppContext) func(next http.Handler) http.Handler {
	cache := appCtx.PC()
	uDAO := appCtx.DM().UsersDAO()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			lgr := appCtx.Lgr()
			ctx := req.Context()
			lgr.Info(fmt.Sprintf("Permission Auth: %+v", permissionName))

			deniedErr := errors.New("Permission denied")

			cookie, err := tools.GetAuthCookie(req)
			if err != nil {
				tools.RequestHttpError(ctx, lgr, res, 403, deniedErr)
				return
			}

			token, id, err := tools.GetAuthParts(cookie)
			if err != nil {
				tools.RequestHttpError(ctx, lgr, res, 403, deniedErr)
				return
			}
			lgr.Info(fmt.Sprintf("User id: %d", id))

			levelID, err := uDAO.GetPrivilegeLevelID(id, token)
			if err != nil {
				tools.RequestHttpError(ctx, lgr, res, 403, deniedErr)
				return
			}

			permitted := cache.HasPermissionByName(levelID, permissionName)
			if !permitted {
				tools.RequestHttpError(ctx, lgr, res, 403, deniedErr)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
