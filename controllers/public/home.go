package public

import (
	"net/http"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
)

type Home struct {
	interfaces.IAppContext
}

func (h *Home) SetAppCtx(ctx interfaces.IAppContext) {
	h.IAppContext = ctx
}

func (r *Home) Path() string {
	return "/"
}

func (hw *Home) PublicRoute(r chi.Router) {
	r.Get("/", hw.redirect_home)
	r.Get("/home", hw.home)
}

func (hw *Home) redirect_home(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/home", http.StatusMovedPermanently)
}

func (hw *Home) home(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	page := pageLayouts.MainPageLayout(pages.Home())
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	page.Render(ctx, res)
}
