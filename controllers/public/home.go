package public

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools/render"
	"github.com/go-chi/chi/v5"
)

type home struct {
	context.AppContext
}

func NewHome(ctx context.AppContext) *home {
	return &home{
		AppContext: ctx,
	}
}

func (r *home) Path() string {
	return "/"
}

func (hw *home) PublicRoute(r chi.Router) {
	r.Get("/", hw.redirect_home)
	r.Get("/home", hw.home)
}

func (hw *home) redirect_home(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/home", http.StatusMovedPermanently)
}

func (hw *home) home(res http.ResponseWriter, req *http.Request) {
	lgr := hw.Lgr("home")
	lgr.Info("Called")
	ctx := req.Context()
	page := pages.Home()
	render.PageMainLayout(req, page).Render(ctx, res)
}
