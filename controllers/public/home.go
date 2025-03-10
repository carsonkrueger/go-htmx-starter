package public

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/templates/layouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Home struct {
	context.WithAppContext
}

func (r *Home) Path() string {
	return "/"
}

func (hw *Home) PublicRoute(r chi.Router) {
	r.Get("/", hw.redirect_home)
	r.Get("/home", hw.home)
	r.Get("/signup", hw.signup)
}

func (hw *Home) redirect_home(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/home", http.StatusMovedPermanently)
}

func (hw *Home) home(res http.ResponseWriter, req *http.Request) {
	lgr := hw.AppCtx.Lgr().With(zap.String("controller", "/home"))
	ctx := req.Context()
	home := layouts.Main(pages.Home())
	err := home.Render(ctx, res)
	if err != nil {
		tools.RequestHttpError(ctx, lgr, res, 500, err)
	}
}

func (hw *Home) signup(res http.ResponseWriter, req *http.Request) {
	lgr := hw.AppCtx.Lgr().With(zap.String("controller", "/signup"))
	ctx := req.Context()
	home := layouts.Main(pages.HomeSignup())
	err := home.Render(ctx, res)
	if err != nil {
		tools.RequestHttpError(ctx, lgr, res, 500, err)
	}
}
