package internal

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/carsonkrueger/main/internal/enums"
	"github.com/carsonkrueger/main/tools"
)

func EnforceAuth(appCtx *AppContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			lgr := appCtx.Lgr
			authToken, err := tools.GetAuthCookie(req)
			if err != nil {
				tools.RequestHttpError(appCtx.Lgr, res, 403, errors.New("Not Authenticated"))
				return
			}
			parts := strings.Split(authToken.Value, "$")
			if len(parts) != 2 {
				tools.RequestHttpError(lgr, res, 403, errors.New("Malformed auth token"))
				return
			}
			auth := parts[0]
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				tools.RequestHttpError(lgr, res, 403, errors.New("Malformed auth token"))
				return
			}
			user, err := appCtx.SM.UsersService.Index(int64(id))
			if err != nil {
				tools.RequestHttpError(lgr, res, 403, errors.New("Malformed auth token"))
				return
			}
			if *user.AuthToken != auth {
				tools.RequestHttpError(lgr, res, 403, errors.New("Malformed auth token"))
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}

func ApplyPermission(p enums.Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			// tokenStr := res.Header().Get("AUTH_TOKEN")
			// parse tokenStr
			// ensure token has permission
			// if not authorized {
			// 		res.WriteHeader(403)
			//		res.Write([]byte(p))
			// }
			next.ServeHTTP(res, req)
		})
	}
}
