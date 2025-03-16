package public

import (
	"net/http"
	"os"
	"path"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type WebPublic struct {
	interfaces.IAppContext
}

func (w *WebPublic) SetAppCtx(ctx interfaces.IAppContext) {
	w.IAppContext = ctx
}

func (wp *WebPublic) Path() string {
	return "/public"
}

func (wp *WebPublic) PublicRoute(r chi.Router) {
	r.Handle("/*", wp.ServePublicDir())
}

func (wp *WebPublic) ServePublicDir() http.Handler {
	lgr := wp.Lgr().With(zap.String("controller", "GET /public"))
	lgr.Info("Initialized WebPublic")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir_path := path.Join(wd, wp.Path())
	handler := http.FileServer(http.Dir(dir_path))
	return http.StripPrefix(wp.Path(), handler)
}
