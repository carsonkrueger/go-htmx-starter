package tools

import (
	"net/http"

	"go.uber.org/zap"
)

func RequestHttpError(logger *zap.Logger, res http.ResponseWriter, err error, code int) {
	logger.Error("HTTP error: ", zap.Int("code", code), zap.Error(err))
	res.WriteHeader(code)
	var e error
	if code >= 500 {
		_, e = res.Write([]byte("An error occurred - Please try again later"))
	} else {
		_, e = res.Write([]byte(err.Error()))
	}
	if e != nil {
		logger.Error("Failed to write error response", zap.Error(err))
	}
}
