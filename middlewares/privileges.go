package middlewares

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/tools"
)

func ApplyPermission(privilegeID int64, appCtx context.AppContext) func(next http.Handler) http.Handler {
	lgr := appCtx.Lgr("MW ApplyPermission")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			levelID := context.GetPrivilegeLevelID(ctx)
			cache := appCtx.SM().PrivilegesService()
			permitted := cache.HasPermissionByID(levelID, privilegeID)
			if !permitted {
				tools.HandleError(req, res, lgr, nil, 403, "Insufficient privileges")
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
