package public

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/templatetargets"
	"github.com/carsonkrueger/main/internal/templates/ui/pages"
	"github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/util"
	"github.com/carsonkrueger/main/pkg/util/render"
	"github.com/carsonkrueger/main/pkg/util/validate"
	"github.com/go-chi/chi/v5"
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

func (s *signUp) postSignup(w http.ResponseWriter, req *http.Request) {
	lgr := s.Lgr("postSignup")
	lgr.Info("Called")
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		common.HandleError(req, w, lgr, err, 403, "Invalid Form")
		return
	}

	form := req.Form
	errs := validate.ValidateSignup(form)
	if len(errs) > 0 {
		common.HandleError(req, w, lgr, errs[0], 400, errs[0].Error())
		return
	}

	salt, _ := util.GenerateSalt()
	authToken, _ := util.GenerateSalt()
	hash := util.HashPassword(form.Get("password"), salt)
	user := model.Users{
		FirstName: form.Get("first_name"),
		LastName:  form.Get("last_name"),
		Email:     form.Get("email"),
		Password:  hash,
		RoleID:    1,
	}

	usersDAO := s.DM().UsersDAO()
	if err := usersDAO.Insert(ctx, &user); err != nil {
		common.HandleError(req, w, lgr, err, 400, "Email already taken")
		return
	}

	session := &model.Sessions{
		UserID: user.ID,
		Token:  authToken,
	}
	sessionDAO := s.DM().SessionsDAO()
	if err := sessionDAO.Insert(ctx, session); err != nil {
		common.HandleError(req, w, lgr, err, 500, "Error creating session")
		return
	}

	util.SetAuthCookie(w, &authToken, constant.AUTH_TOKEN_KEY)

	render.Layout(ctx, req, w, templatetargets.Main, pages.Login())
}

func (s *signUp) getSignup(w http.ResponseWriter, req *http.Request) {
	lgr := s.Lgr("getSignup")
	lgr.Info("Called")
	ctx := req.Context()
	render.Layout(ctx, req, w, templatetargets.Main, pages.Signup())
}
