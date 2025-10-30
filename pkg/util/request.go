package util

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func SetAuthCookie(w http.ResponseWriter, authToken *string, key string) {
	cookie := http.Cookie{
		Name:     key,
		Value:    *authToken,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

func GetAuthCookie(req *http.Request, key string) (*http.Cookie, error) {
	cookie, err := req.Cookie(key)
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

func HandleRedirect(req *http.Request, w http.ResponseWriter, url string) {
	if IsHxRequest(req) {
		w.Header().Set("Hx-Redirect", url)
	} else {
		http.Redirect(w, req, url, http.StatusFound)
	}
}
