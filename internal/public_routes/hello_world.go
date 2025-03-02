package public_routes

import (
	"context"
	"net/http"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/templates/layouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/go-chi/chi/v5"
)

type HelloWorld struct {
	types.WithAppContext
}

func (r *HelloWorld) Path() string {
	return "/"
}

func (hw *HelloWorld) PublicRoute(r chi.Router) {
	r.Get("/", hw.redirect_home)
	r.Get("/home", hw.home)
	r.Get("/about", hw.about)
}

func (hw *HelloWorld) redirect_home(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/home", http.StatusMovedPermanently)
}

func (hw *HelloWorld) home(res http.ResponseWriter, req *http.Request) {
	hw.AppCtx.Lgr.Info("Logging HelloWorld route")
	home := layouts.Main(pages.Home(), pages.Home())
	home.Render(context.Background(), res)
}

func (hw *HelloWorld) about(res http.ResponseWriter, req *http.Request) {
	// err := hw.GetCtx().Templates.Render(res, "about.html", nil)
	// types.ReportIfErr(err, nil)
}
