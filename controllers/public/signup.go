package public

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/validate"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type signUp struct {
	context.AppContext
}

func NewSignUp(ctx context.AppContext) *signUp {
	return &signUp{AppContext: ctx}
}

func (s *signUp) Path() string {
	return "/signup"
}

func (s *signUp) PublicRoute(r chi.Router) {
	r.Get("/", s.getSignup)
	r.Post("/", s.postSignup)
}

func (s *signUp) postSignup(res http.ResponseWriter, req *http.Request) {
	lgr := s.Lgr("postSignup")
	lgr.Info("Called")
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		tools.HandleError(req, res, lgr, err, 403, "Invalid Form")
		return
	}

	form := req.Form
	errs := validate.ValidateSignup(form)
	if len(errs) > 0 {
		lgr.Warn("Validation errors", zap.Errors("Signup Form", errs))
		res.WriteHeader(422)
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
		PrivilegeLevelID: 1000,
	}

	usersDAO := s.DM().UsersDAO()
	if err := usersDAO.Insert(&user, s.DB()); err != nil {
		lgr.Warn("Could not insert user", zap.Error(err))
		res.WriteHeader(422)
		noti := datadisplay.AddTextToast(datadisplay.Warning, "Email taken", 0)
		noti.Render(ctx, res)
		return
	}

	session := &model.Sessions{
		UserID: user.ID,
		Token:  authToken,
	}
	sessionDAO := s.DM().SessionsDAO()
	if err := sessionDAO.Insert(session, s.DB()); err != nil {
		lgr.Error("Could not insert session", zap.Error(err))
		res.WriteHeader(500)
		noti := datadisplay.AddTextToast(datadisplay.Error, "Error creating session", 0)
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

func (s *signUp) getSignup(res http.ResponseWriter, req *http.Request) {
	lgr := s.Lgr("getSignup")
	lgr.Info("Called")
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	page := pages.Signup()
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(pageLayouts.MainPageLayout(page))
	}
	page.Render(ctx, res)
}
