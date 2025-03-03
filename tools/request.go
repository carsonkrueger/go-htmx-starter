package tools

import (
	"net/http"

	"go.uber.org/zap"
)

func RequestHttpError(logger *zap.Logger, res http.ResponseWriter, code int, errs ...error) {
	res.WriteHeader(code)
	for _, e := range errs {
		logger.Error("HTTP error: ", zap.Int("code", code), zap.Error(e))
		var e error
		if code >= 500 {
			_, e = res.Write([]byte("An error occurred - Please try again later"))
			return
		} else {
			_, e = res.Write([]byte(e.Error()))
		}
		if e != nil {
			logger.Error("Failed to write error response", zap.Error(e))
		}
	}
}
