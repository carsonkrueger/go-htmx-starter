package tools

import (
	"net/http"

	"go.uber.org/zap"
)

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
