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
	lgr := appCtx.Lgr("MW ApplyPermission")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			lgr.Info(fmt.Sprintf("Permission Auth: %+v", permissionName))

			deniedErr := errors.New("Permission denied")

			levelID := context.GetPrivilegeLevelID(ctx)
			cache := appCtx.SM().PrivilegesService()
			permitted := cache.HasPermissionByName(levelID, permissionName)
			if !permitted {
				tools.RequestHttpError(ctx, lgr, res, 403, deniedErr, "Permission denied")
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
