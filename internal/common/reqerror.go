package common

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/templates/ui/partials/toast"
	tuitoast "github.com/carsonkrueger/main/pkg/templui/toast"
	"github.com/carsonkrueger/main/pkg/util"
	"go.uber.org/zap"
)

func HandleError(req *http.Request, w http.ResponseWriter, lgr *zap.Logger, err error, status int, msg string) {
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
	if util.IsHxRequest(req) {
		w.Header().Add("hx-toast-err", "true")
		w.WriteHeader(status)
		toast.AddToastErrors(tuitoast.Props{Description: msg}).Render(ctx, w)
	} else {
		w.WriteHeader(status)
		w.Write([]byte(msg))
	}
}
