package public_routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/validate"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Auth struct {
	types.WithAppContext
}

func (a *Auth) Path() string {
	return "/auth"
}

func (a *Auth) PublicRoute(r chi.Router) {
	r.Post("/login", a.login)
	r.Post("/signup", a.signup)
}

func (a *Auth) login(res http.ResponseWriter, req *http.Request) {
	lgr := a.AppCtx.Lgr.With(zap.String("controller", "/auth/login"))

	if err := req.ParseForm(); err != nil {
		tools.RequestHttpError(lgr, res, 400, errors.New("Error parsing form"))
		return
	}

	form := req.Form
	errs := validate.ValidateLogin(form)

	if len(errs) > 0 {
		tools.RequestHttpError(lgr, res, 400, errs...)
		return
	}

	email := form.Get("email")
	password := form.Get("password")

	usersService := a.AppCtx.SM.UsersService
	user, err := usersService.GetByEmail(email)
	if err != nil {
		tools.RequestHttpError(lgr, res, 403, errors.New("Invalid email or password"))
		return
	}

	parts := strings.Split(user.Password, "$")
	if len(parts) != 2 {
		tools.RequestHttpError(lgr, res, 500, errors.New(fmt.Sprintf("Invalid hash: %d", user.ID)))
		return
	}

	hash := tools.HashPassword(password, parts[0])
	if user.Password != hash {
		tools.RequestHttpError(a.AppCtx.Lgr, res, 403, errors.New("Invalid email or password"))
		return
	}

	res.Write(fmt.Appendf(nil, "%s - %s", email, password))
}

func (a *Auth) signup(res http.ResponseWriter, req *http.Request) {
	lgr := a.AppCtx.Lgr.With(zap.String("controller", "/signup/"))

	if err := req.ParseForm(); err != nil {
		tools.RequestHttpError(lgr, res, 400, errors.New("Error parsing form"))
		return
	}

	form := req.Form
	errs := validate.ValidateSignup(form)
	if len(errs) > 0 {
		tools.RequestHttpError(lgr, res, 400, errs...)
		return
	}

	salt, _ := tools.GenerateSalt()
	auth_token, _ := tools.GenerateSalt()
	hash := tools.HashPassword(form.Get("password"), salt)
	user := model.Users{
		FirstName: form.Get("first_name"),
		LastName:  form.Get("last_name"),
		Email:     form.Get("email"),
		Password:  hash,
		AuthToken: auth_token,
	}

	_, err := a.AppCtx.SM.UsersService.Insert(&user)
	if err != nil {
		tools.RequestHttpError(lgr, res, 500, err)
		return
	}

	cookie := http.Cookie{
		Name:     "ghx_auth_token",
		Value:    auth_token,
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Hx-Redirect", "/home")
	res.Write([]byte{})
}
