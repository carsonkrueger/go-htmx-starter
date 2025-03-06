package tools

import (
	"net/http"

	"go.uber.org/zap"
)

var AUTH_TOKEN_KEY string = "ghx_auth_token"

func RequestHttpError(logger *zap.Logger, res http.ResponseWriter, code int, errs ...error) {
	res.WriteHeader(code)
	for _, e := range errs {
		if code >= 500 {
			logger.Error("Internal error: ", zap.Int("status code", code), zap.Error(e))
			res.Write([]byte("An error occurred - Please try again later"))
		} else {
			logger.Info("Request error: ", zap.Int("status code", code), zap.Error(e))
			res.Write([]byte(e.Error()))
		}
	}
}

func SetAuthCookie(res http.ResponseWriter, authToken *string) {
	cookie := http.Cookie{
		Name:     AUTH_TOKEN_KEY,
		Value:    *authToken,
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
}

func GetAuthCookie(req *http.Request) (*http.Cookie, error) {
	cookie, err := req.Cookie(AUTH_TOKEN_KEY)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}
