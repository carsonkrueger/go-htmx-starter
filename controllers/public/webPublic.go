package public

import (
	"net/http"
	"os"
	"path"

	"github.com/carsonkrueger/main/context"
	"github.com/go-chi/chi/v5"
)

type webPublic struct {
	context.AppContext
}

func NewWebPublic(ctx context.AppContext) *webPublic {
	return &webPublic{
		AppContext: ctx,
	}
}

func (wp *webPublic) Path() string {
	return "/public"
}

func (wp *webPublic) PublicRoute(r chi.Router) {
	r.Handle("/*", wp.ServePublicDir())
}

func (wp *webPublic) ServePublicDir() http.Handler {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir_path := path.Join(wd, wp.Path())
	handler := http.FileServer(http.Dir(dir_path))
	return http.StripPrefix(wp.Path(), handler)
}
