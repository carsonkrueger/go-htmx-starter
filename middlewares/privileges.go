package middlewares

import (
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

			levelID := context.GetPrivilegeLevelID(ctx)
			cache := appCtx.SM().PrivilegesService()
			permitted := cache.HasPermissionByName(levelID, permissionName)
			if !permitted {
				tools.HandleError(req, res, lgr, nil, 403, "Malformed auth token")
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
