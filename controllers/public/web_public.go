package public

import (
	"net/http"
	"os"
	"path"

	"github.com/carsonkrueger/main/context"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type WebPublic struct {
	context.WithAppContext
}

func (wp *WebPublic) Path() string {
	return "/public"
}

func (wp *WebPublic) PublicRoute(r chi.Router) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	r.Handle("/*", wp.ServePublicDir(wd))
}

func (wp *WebPublic) ServePublicDir(wd string) http.Handler {
	lgr := wp.AppCtx.Lgr().With(zap.String("controller", "GET /public"))
	lgr.Info("Initialized WebPublic")
	dir_path := path.Join(wd, wp.Path())
	handler := http.FileServer(http.Dir(dir_path))
	return http.StripPrefix(wp.Path(), handler)
}
