package public

import (
	"net/http"
	"os"
	"path"

	"github.com/carsonkrueger/main/context"
	"github.com/go-chi/chi/v5"
)

type WebPublic struct {
	context.WithAppContext
}

func (w *WebPublic) Path() string {
	return "/public"
}

func (w *WebPublic) PublicRoute(r chi.Router) {
	r.Handle("/*", w.ServePublicDir())
}

func (w *WebPublic) ServePublicDir() http.Handler {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir_path := path.Join(wd, w.Path())
	handler := http.FileServer(http.Dir(dir_path))
	return http.StripPrefix(w.Path(), handler)
}
