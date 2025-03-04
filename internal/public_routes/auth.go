package public_routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/tools"
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
	form_errs := []error{}

	email := form.Get("email")
	password := form.Get("password")
	if email == "" {
		form_errs = append(form_errs, errors.New("Missing email"))
	}
	if password == "" {
		form_errs = append(form_errs, errors.New("Missing password"))
	}
	if len(form_errs) > 0 {
		tools.RequestHttpError(lgr, res, 400, form_errs...)
		return
	}

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
	lgr := a.AppCtx.Lgr.With(zap.String("Route", "/auth/signup"))

	if err := req.ParseForm(); err != nil {
		tools.RequestHttpError(lgr, res, 400, errors.New("Error parsing form"))
		return
	}

	form := req.Form
	form_errs := []error{}

	email := form.Get("email")
	password := form.Get("password")
	first := form.Get("first_name")
	last := form.Get("last_name")
	if email == "" {
		form_errs = append(form_errs, errors.New("Missing email"))
	}
	if password == "" {
		form_errs = append(form_errs, errors.New("Missing password"))
	}
	if first == "" {
		form_errs = append(form_errs, errors.New("Missing first name"))
	}
	if last == "" {
		form_errs = append(form_errs, errors.New("Missing last name"))
	}
	if len(form_errs) > 0 {
		tools.RequestHttpError(lgr, res, 400, form_errs...)
		return
	}

	salt, _ := tools.GenerateSalt()
	hash := tools.HashPassword(password, salt)

	usersService := a.AppCtx.SM.UsersService
	_, err := usersService.Insert(model.Users{
		Email:     email,
		Password:  hash,
		FirstName: first,
		LastName:  last,
	})
	if err != nil {
		tools.RequestHttpError(lgr, res, 403, errors.New("Failed to create user"))
		return
	}

	// Generate JWT token
	// token, err := tools.GenerateJWT(user.ID)
	// if err != nil {
	// 	tools.RequestHttpError(lgr, res, 500, errors.New("Failed to generate JWT token"))
	// 	return
	// }

	res.Write(fmt.Appendf(nil, "%s - %s", email, password))
}
