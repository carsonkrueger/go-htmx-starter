package public

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/validate"
	"github.com/go-chi/chi/v5"
)

type SignUp struct {
	context.WithAppContext
}

func (s *SignUp) Path() string {
	return "/signup"
}

func (s *SignUp) PublicRoute(r chi.Router) {
	r.Get("/", s.getSignup)
	r.Post("/", s.postSignup)
}

func (s *SignUp) postSignup(res http.ResponseWriter, req *http.Request) {
	// lgr := s.AppCtx.Lgr().With(zap.String("controller", "/auth/signup"))
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		noti := datadisplay.AddTextToast(models.Error, "Error parsing form", 0)
		noti.Render(ctx, res)
		return
	}

	form := req.Form
	errs := validate.ValidateSignup(form)
	if len(errs) > 0 {
		noti := datadisplay.AddToastErrors(0, errs...)
		noti.Render(ctx, res)
		return
	}

	salt, _ := tools.GenerateSalt()
	authToken, _ := tools.GenerateSalt()
	hash := tools.HashPassword(form.Get("password"), salt)
	user := model.Users{
		FirstName:        form.Get("first_name"),
		LastName:         form.Get("last_name"),
		Email:            form.Get("email"),
		Password:         hash,
		AuthToken:        &authToken,
		PrivilegeLevelID: 1000,
	}

	dao := s.AppCtx.DM().UsersDAO()
	_, err := dao.Insert(&user)
	if err != nil {
		res.WriteHeader(422)
		noti := datadisplay.AddToastErrors(0, err)
		noti.Render(ctx, res)
		return
	}

	tools.SetAuthCookie(res, &authToken)

	hxRequest := tools.IsHxRequest(req)
	page := pages.Login()
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(pageLayouts.MainPageLayout(page))
	}
	page.Render(ctx, res)
}

func (s *SignUp) getSignup(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	page := pages.Signup()
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(pageLayouts.MainPageLayout(page))
	}
	page.Render(ctx, res)
}
