package public

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/validate"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Auth struct {
	context.WithAppContext
}

func (a *Auth) Path() string {
	return "/auth"
}

func (a *Auth) PublicRoute(r chi.Router) {
	r.Post("/login", a.login)
	r.Post("/signup", a.signup)
}

func (a *Auth) login(res http.ResponseWriter, req *http.Request) {
	lgr := a.AppCtx.Lgr().With(zap.String("controller", "/auth/login"))
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		tools.RequestHttpError(ctx, lgr, res, 400, errors.New("Error parsing form"))
		return
	}

	form := req.Form
	errs := validate.ValidateLogin(form)

	if len(errs) > 0 {
		tools.RequestHttpError(ctx, lgr, res, 422, errs...)
		return
	}

	email := form.Get("email")
	password := form.Get("password")

	usersService := a.AppCtx.SM().UsersService()
	authToken, err := usersService.Login(email, password)
	if err != nil {
		lgr.Info("Error logging in user", zap.Error(err))
		tools.RequestHttpError(ctx, lgr, res, 422, errors.New("Invalid username or password"))
		return
	}

	tools.SetAuthCookie(res, authToken)
	fmt.Println(*authToken)
}

func (a *Auth) signup(res http.ResponseWriter, req *http.Request) {
	lgr := a.AppCtx.Lgr().With(zap.String("controller", "/signup/"))
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		tools.RequestHttpError(ctx, lgr, res, 400, errors.New("Error parsing form"))
		return
	}

	form := req.Form
	errs := validate.ValidateSignup(form)
	if len(errs) > 0 {
		tools.RequestHttpError(ctx, lgr, res, 400, errs...)
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

	dao := a.AppCtx.DM().UsersDAO()
	_, err := dao.Insert(&user)
	if err != nil {
		tools.RequestHttpError(ctx, lgr, res, 500, err)
		return
	}

	tools.SetAuthCookie(res, &authToken)
	res.Header().Set("HX-Redirect", "/home")
}
