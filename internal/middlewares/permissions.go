package middlewares

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/enums"
)

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
