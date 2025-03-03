package public_routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
)

type Auth struct {
	types.WithAppContext
}

func (a *Auth) Path() string {
	return "/auth"
}

func (a *Auth) PublicRoute(r chi.Router) {
	r.Post("/login", a.login)
}

func (a *Auth) login(res http.ResponseWriter, req *http.Request) {
	a.AppCtx.Lgr.Info("Auth login route")

	if err := req.ParseForm(); err != nil {
		tools.RequestHttpError(a.AppCtx.Lgr, res, 400, errors.New("Error parsing form"))
		return
	}

	form := req.Form
	form_errs := []error{}

	email := form.Get("email")
	if email == "" {
		form_errs = append(form_errs, errors.New("Missing email"))
	}

	password := form.Get("password")
	if password == "" {
		form_errs = append(form_errs, errors.New("Missing password"))
	}

	if len(form_errs) > 0 {
		tools.RequestHttpError(a.AppCtx.Lgr, res, 400, form_errs...)
		return
	}

	res.Write(fmt.Appendf(nil, "%s - %s", email, password))
}
