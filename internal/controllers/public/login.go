package public

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/controllers/private"
	"github.com/carsonkrueger/main/internal/templates/templatetargets"
	"github.com/carsonkrueger/main/internal/templates/ui/pages"
	"github.com/carsonkrueger/main/internal/templates/ui/tables"
	"github.com/carsonkrueger/main/pkg/util"
	"github.com/carsonkrueger/main/pkg/util/render"
	"github.com/carsonkrueger/main/pkg/util/validate"
	"github.com/go-chi/chi/v5"
)

type login struct {
	*context.AppContext
}

func NewLogin(ctx *context.AppContext) *login {
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

func (l *login) postLogin(w http.ResponseWriter, req *http.Request) {
	lgr := l.Lgr("postLogin")
	lgr.Info("Called")
	ctx := req.Context()
	db := context.GetDB(ctx)
	db.Query("select 1;")

	if err := req.ParseForm(); err != nil {
		common.HandleError(req, w, lgr, nil, 400, "Error parsing form")
		return
	}

	form := req.Form
	errs := validate.ValidateLogin(form)

	if len(errs) > 0 {
		common.HandleError(req, w, lgr, errs[0], 400, errs[0].Error())
		return
	}

	email := form.Get("email")
	password := form.Get("password")

	usersService := l.SM().UsersService()
	authToken, err := usersService.Login(ctx, email, password, req)
	if err != nil {
		common.HandleError(req, w, lgr, err, 401, "Invalid username or password")
		return
	}

	util.SetAuthCookie(w, authToken, constant.AUTH_TOKEN_KEY)

	hxRequest := util.IsHxRequest(req)
	if hxRequest {
		dao := l.DM().UsersDAO()
		users, err := dao.GetUserPrivilegeJoinAll(ctx)
		if err != nil || users == nil {
			common.HandleError(req, w, lgr, err, 500, "Error fetching privileges")
			return
		}

		if len(users) == 0 {
			return
		}

		allRoles, err := l.DM().RolesDAO().GetAll(ctx)
		if err != nil || allRoles == nil {
			common.HandleError(req, w, lgr, err, 500, "Error fetching roles")
			return
		}

		table := tables.ManageUsersTable(users, allRoles)
		render.Tab(req, private.UserManagementTabModels, 0, table).Render(ctx, w)
	}
}

func (l *login) getLogin(w http.ResponseWriter, req *http.Request) {
	lgr := l.Lgr("getLogin")
	lgr.Info("Called")
	ctx := req.Context()
	render.Layout(ctx, req, w, templatetargets.Main, pages.Login())
}
