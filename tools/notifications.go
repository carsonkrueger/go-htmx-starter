package tools

import (
	"net/http"

	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
)

func RenderErrorNotification(req *http.Request, res http.ResponseWriter, msg string, duration int) {
	if IsHxRequest(req) {
		ctx := req.Context()
		datadisplay.AddTextToast(models.Error, msg, duration).Render(ctx, res)
	}
}
