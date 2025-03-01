package public_routes

import (
	"net/http"
	"os"
	"path"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/go-chi/chi/v5"
)

type WebPublic struct {
	types.WithAppContext
}

func (w *WebPublic) Path() string {
	return "/public"
}

func (w *WebPublic) PublicRoute(r chi.Router) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir_path := path.Join(wd, w.Path())
	dir := http.FileServer(http.Dir(dir_path))
	r.Handle("/*", http.StripPrefix(w.Path(), dir))
}
