package tools

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/carsonkrueger/main/constant"
	"go.uber.org/zap"
)

func RequestHttpError(ctx context.Context, logger *zap.Logger, res http.ResponseWriter, code int, errs ...error) {
	res.WriteHeader(code)
	for _, e := range errs {
		if code >= 500 {
			logger.Error("Internal error: ", zap.Int("status code", code), zap.Error(e))
			res.Write([]byte("An error occurred - Please try again later"))
		} else {
			logger.Info("Request error: ", zap.Int("status code", code), zap.Error(e))
		}
	}
}

func SetAuthCookie(res http.ResponseWriter, authToken *string) {
	cookie := http.Cookie{
		Name:     constant.AUTH_TOKEN_KEY,
		Value:    *authToken,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(res, &cookie)
}

func GetAuthCookie(req *http.Request) (*http.Cookie, error) {
	cookie, err := req.Cookie(constant.AUTH_TOKEN_KEY)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func GetAuthParts(cookie *http.Cookie) (string, int64, error) {
	if cookie == nil {
		return "", 0, errors.New("cookie is nil")
	}
	parts := strings.Split(cookie.Value, "$")
	if parts == nil || len(parts) != 2 {
		return "", 0, errors.New("invalid cookie")
	}
	token := parts[0]
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, errors.New("Could not parse user ID")
	}
	return token, int64(id), err
}

func IsHxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
