package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/tools"
)

func ApplyPermission(permissionName string, appCtx interfaces.IAppContext) func(next http.Handler) http.Handler {
	cache := appCtx.PC()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			lgr := appCtx.Lgr()
			ctx := req.Context()
			lgr.Info(fmt.Sprintf("Permission Auth: %+v", permissionName))

			deniedErr := errors.New("Permission denied")

			levelID := context.GetPrivilegeID(ctx)
			permitted := cache.HasPermissionByName(levelID, permissionName)
			if !permitted {
				tools.RequestHttpError(ctx, lgr, res, 403, deniedErr)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
