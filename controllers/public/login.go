package public

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/controllers/private"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/render"
	"github.com/carsonkrueger/main/tools/validate"
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

	if err := req.ParseForm(); err != nil {
		lgr.Error("Could not parse form", zap.Error(err))
		res.WriteHeader(422)
		noti := datadisplay.AddTextToast(datadisplay.Error, "Error parsing form", 0)
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
		noti := datadisplay.AddTextToast(datadisplay.Error, "Invalid username or password", 5)
		noti.Render(ctx, res)
		return
	}

	tools.SetAuthCookie(res, authToken)

	hxRequest := tools.IsHxRequest(req)
	if hxRequest {
		dao := l.DM().UsersDAO()
		users, err := dao.GetUserPrivilegeJoinAll()
		if err != nil || users == nil {
			tools.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
			return
		}

		if len(*users) == 0 {
			datadisplay.AddTextToast(datadisplay.Warning, "No Users Found", 5).Render(ctx, res)
			return
		}

		allLevels, err := l.DM().PrivilegeLevelsDAO().Index(nil, l.DB())
		if err != nil || allLevels == nil {
			tools.HandleError(req, res, lgr, err, 500, "Error fetching privilege levels")
			return
		}

		rows := l.SM().PrivilegesService().UserPrivilegeLevelJoinAsRowData(*users, allLevels)
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
