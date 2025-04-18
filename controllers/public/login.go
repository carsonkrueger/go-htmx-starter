package public

import (
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/templates/partials"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/validate"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type login struct {
	interfaces.IAppContext
}

func NewLogin(ctx interfaces.IAppContext) *login {
	return &login{
		IAppContext: ctx,
	}
}

func (l *login) Path() string {
	return "/login"
}

func (l *login) PublicRoute(r chi.Router) {
	r.Get("/", l.getLogin)
	r.Post("/", l.postLogin)
}

func (l *login) postLogin(res http.ResponseWriter, req *http.Request) {
	lgr := l.Lgr("postLogin")
	lgr.Info("Called")
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		lgr.Error("Could not parse form", zap.Error(err))
		res.WriteHeader(422)
		noti := datadisplay.AddTextToast(models.Error, "Error parsing form", 0)
		noti.Render(ctx, res)
		return
	}

	form := req.Form
	errs := validate.ValidateLogin(form)

	if len(errs) > 0 {
		lgr.Warn("Validation errors", zap.Errors("Login Form", errs))
		res.WriteHeader(422)
		noti := datadisplay.AddToastErrors(0, errs...)
		noti.Render(ctx, res)
		return
	}

	email := form.Get("email")
	password := form.Get("password")

	usersService := l.SM().UsersService()
	authToken, err := usersService.Login(email, password, req)
	if err != nil {
		lgr.Warn("Could not login")
		res.WriteHeader(422)
		noti := datadisplay.AddTextToast(models.Error, "Invalid username or password", 5)
		noti.Render(ctx, res)
		return
	}

	tools.SetAuthCookie(res, authToken)

	hxRequest := tools.IsHxRequest(req)
	if hxRequest {
		content := partials.Redirect("/user_management/users", "#"+pageLayouts.MainContentID, builders.GET, "true")
		content.Render(ctx, res)
	}
}

func (l *login) getLogin(res http.ResponseWriter, req *http.Request) {
	lgr := l.Lgr("getLogin")
	lgr.Info("Called")
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	page := pageLayouts.MainPageLayout(pages.Login())
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	page.Render(ctx, res)
}
