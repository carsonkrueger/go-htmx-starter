package tools

import (
	"net/http"

	"github.com/carsonkrueger/main/templates/datadisplay"
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
		datadisplay.AddTextToast(datadisplay.Error, msg, 0).Render(ctx, res)
	} else {
		res.Write([]byte(msg))
		res.WriteHeader(status)
	}
}
