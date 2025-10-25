package public

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/datadisplay"
	"github.com/carsonkrueger/main/internal/templates/page_layouts"
	"github.com/carsonkrueger/main/internal/templates/pages"
	"github.com/carsonkrueger/main/pkg/model/db/auth"
	"github.com/carsonkrueger/main/pkg/util"
	"github.com/carsonkrueger/main/pkg/util/slice"
	"github.com/carsonkrueger/main/pkg/util/validate"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type signUp struct {
	*context.AppContext
}

func NewSignUp(ctx *context.AppContext) *signUp {
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
		util.HandleError(req, res, lgr, err, 403, "Invalid Form")
		return
	}

	form := req.Form
	errs := validate.ValidateSignup(form)
	if len(errs) > 0 {
		lgr.Warn("Validation errors", zap.Errors("Signup Form", errs))
		res.WriteHeader(422)
		datadisplay.AddToastErrors(slice.Map(errs, error.Error)...).Render(ctx, res)
		return
	}

	salt, _ := util.GenerateSalt()
	authToken, _ := util.GenerateSalt()
	hash := util.HashPassword(form.Get("password"), salt)
	user := auth.Users{
		FirstName: form.Get("first_name"),
		LastName:  form.Get("last_name"),
		Email:     form.Get("email"),
		Password:  hash,
		RoleID:    1,
	}

	usersDAO := s.DM().UsersDAO()
	if err := usersDAO.Insert(ctx, &user); err != nil {
		lgr.Warn("Could not insert user", zap.Error(err))
		res.WriteHeader(422)
		datadisplay.AddToastErrors("Email taken").Render(ctx, res)
		return
	}

	session := &auth.Sessions{
		UserID: user.ID,
		Token:  authToken,
	}
	sessionDAO := s.DM().SessionsDAO()
	if err := sessionDAO.Insert(ctx, session); err != nil {
		lgr.Error("Could not insert session", zap.Error(err))
		res.WriteHeader(500)
		datadisplay.AddToastErrors("Error creating session").Render(ctx, res)
		return
	}

	util.SetAuthCookie(res, &authToken, constant.AUTH_TOKEN_KEY)

	hxRequest := util.IsHxRequest(req)
	page := pages.Login()
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = page_layouts.Index(page_layouts.MainPageLayout(page))
	}
	page.Render(ctx, res)
}

func (s *signUp) getSignup(res http.ResponseWriter, req *http.Request) {
	lgr := s.Lgr("getSignup")
	lgr.Info("Called")
	ctx := req.Context()
	hxRequest := util.IsHxRequest(req)
	page := pages.Signup()
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = page_layouts.Index(page_layouts.MainPageLayout(page))
	}
	page.Render(ctx, res)
}
