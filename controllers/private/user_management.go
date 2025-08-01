package private

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/page_layouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/render"
)

const (
	UserManagementTabsGet   = "UserManagementTabsGet"
	UserManagementUsersGet  = "UserManagementUsersGet"
	UserManagementLevelsGet = "UserManagementLevelsGet"
)

var UserManagementTabModels = []page_layouts.TabModel{
	{Title: "Users", PushUrl: true, HxGet: "/user_management/users"},
	{Title: "Privilege Levels", PushUrl: true, HxGet: "/user_management/levels"},
}

type userManagement struct {
	context.AppContext
}

func NewUserManagement(ctx context.AppContext) *userManagement {
	return &userManagement{
		AppContext: ctx,
	}
}

func (um userManagement) Path() string {
	return "/user_management"
}

func (um *userManagement) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	// b.NewHandle().Register(builders.GET, "/tabs", um.userManagementTabsGet).SetPermissionName(UserManagementTabsGet).Build()
	b.NewHandler().Register(builders.GET, "/users", um.userManagementUsersGet).SetPermissionName(UserManagementUsersGet).Build(ctx)
	b.NewHandler().Register(builders.GET, "/levels", um.userManagementLevelsGet).SetPermissionName(UserManagementLevelsGet).Build(ctx)
}

func (um *userManagement) userManagementUsersGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementUsersGet")
	lgr.Info("Called")
	ctx := req.Context()

	dao := um.DM().UsersDAO()
	users, err := dao.GetUserPrivilegeJoinAll(ctx)
	if err != nil || users == nil {
		tools.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}

	if len(*users) == 0 {
		datadisplay.AddTextToast(datadisplay.Warning, "No Users Found", 5).Render(ctx, res)
		return
	}

	allLevels, err := um.DM().PrivilegeLevelsDAO().Index(ctx, nil)
	if err != nil || allLevels == nil {
		tools.HandleError(req, res, lgr, err, 500, "Error fetching privilege levels")
		return
	}

	rows := um.SM().PrivilegesService().UserPrivilegeLevelJoinAsRowData(ctx, *users, allLevels)
	page := pages.UserManagementUsers(rows)
	render.Tab(req, UserManagementTabModels, 0, page).Render(ctx, res)
}

func (um *userManagement) userManagementLevelsGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementLevelsGet")
	lgr.Info("Called")
	ctx := req.Context()

	privileges, err := um.DM().PrivilegeDAO().GetAllJoined(ctx)
	if err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}
	rows := um.SM().PrivilegesService().JoinedPrivilegesAsRowData(ctx, privileges)

	page := pages.UserManagementLevels(rows)
	render.Tab(req, UserManagementTabModels, 1, page).Render(ctx, res)
}
