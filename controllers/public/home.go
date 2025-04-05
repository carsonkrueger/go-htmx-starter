package public

import (
	"net/http"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
)

type home struct {
	interfaces.IAppContext
}

func NewHome(ctx interfaces.IAppContext) *home {
	return &home{
		IAppContext: ctx,
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
	hxRequest := tools.IsHxRequest(req)
	page := pageLayouts.MainPageLayout(pages.Home())
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	page.Render(ctx, res)
}
