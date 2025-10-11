package public

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/controllers/private"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/util"
	"github.com/carsonkrueger/main/util/render"
	"github.com/carsonkrueger/main/util/slice"
	"github.com/carsonkrueger/main/util/validate"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type login struct {
	context.AppContext
}

func NewLogin(ctx context.AppContext) *login {
	return &login{
		AppContext: ctx,
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
	db := context.GetDB(ctx)
	db.Query("select 1;")

	if err := req.ParseForm(); err != nil {
		// lgr.Error("Could not parse form", zap.Error(err))
		// res.WriteHeader(422)
		// datadisplay.AddToastErrors("Error parsing form").Render(ctx, res)
		util.HandleError(req, res, lgr, nil, 400, "Error parsing form")
		return
	}

	form := req.Form
	errs := validate.ValidateLogin(form)

	if len(errs) > 0 {
		lgr.Warn("Validation errors", zap.Errors("Login Form", errs))
		res.WriteHeader(422)
		datadisplay.AddToastErrors(slice.Map(errs, error.Error)...).Render(ctx, res)
		return
	}

	email := form.Get("email")
	password := form.Get("password")

	usersService := l.SM().UsersService()
	authToken, err := usersService.Login(ctx, email, password, req)
	if err != nil {
		util.HandleError(req, res, lgr, err, 401, "Invalid username or password")
		return
	}

	util.SetAuthCookie(res, authToken)

	hxRequest := util.IsHxRequest(req)
	if hxRequest {
		dao := l.DM().UsersDAO()
		users, err := dao.GetUserPrivilegeJoinAll(ctx)
		if err != nil || users == nil {
			util.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
			return
		}

		if len(*users) == 0 {
			return
		}

		allRoles, err := l.DM().RolesDAO().GetAll(ctx)
		if err != nil || allRoles == nil {
			datadisplay.AddToastErrors("Error fetching roles").Render(ctx, res)
			return
		}

		rows := l.SM().PrivilegesService().UserRoleJoinAsRowData(ctx, *users, allRoles)
		page := pages.UserManagementUsers(rows)
		render.Tab(req, private.UserManagementTabModels, 0, page).Render(ctx, res)
	}
}

func (l *login) getLogin(res http.ResponseWriter, req *http.Request) {
	lgr := l.Lgr("getLogin")
	lgr.Info("Called")
	ctx := req.Context()
	page := pages.Login()
	render.PageMainLayout(req, page).Render(ctx, res)
}
