package util

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/templates/datadisplay"
	"go.uber.org/zap"
)

func HandleError(req *http.Request, res http.ResponseWriter, lgr *zap.Logger, err error, status int, msg string) {
	if status < 400 {
		return
	}
	ctx := req.Context()
	if status >= 400 && status < 500 {
		lgr.Warn(msg, zap.Error(err))
	} else if status >= 500 {
		lgr.Error(msg, zap.Error(err))
		msg = "Internal Server Error"
	}
	if IsHxRequest(req) {
		res.Header().Add("hx-toast-err", "true")
		res.WriteHeader(status)
		datadisplay.AddToastErrors(msg).Render(ctx, res)
	} else {
		res.WriteHeader(status)
		res.Write([]byte(msg))
	}
}
