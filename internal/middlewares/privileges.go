package middlewares

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/context"
)

func ApplyPrivileges(privileges []int64, appCtx *context.AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			lgr := context.GetLogger(ctx, "mw.ApplyPrivileges")

			roleID := context.GetRoleID(ctx)
			privelegeService := appCtx.SM().PrivilegesService()
			permitted := privelegeService.HasPermissionsByIDS(ctx, roleID, privileges)
			if !permitted {
				common.HandleError(req, w, lgr, nil, 403, "Insufficient privileges")
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
