package private

import (
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
	"github.com/carsonkrueger/main/tools/render"
)

const (
	UserManagementTabsGet   = "UserManagementTabsGet"
	UserManagementUsersGet  = "UserManagementUsersGet"
	UserManagementLevelsGet = "UserManagementLevelsGet"
)

var tabModels = []pageLayouts.TabModel{
	{Title: "Users", PushUrl: true, HxGet: "/user_management/users"},
	{Title: "Privilege Levels", PushUrl: true, HxGet: "/user_management/levels"},
}

type userManagement struct {
	interfaces.IAppContext
}

func NewUserManagement(ctx interfaces.IAppContext) *userManagement {
	return &userManagement{
		IAppContext: ctx,
	}
}

func (um userManagement) Path() string {
	return "/user_management"
}

func (um *userManagement) PrivateRoute(b *builders.PrivateRouteBuilder) {
	// b.NewHandle().Register(builders.GET, "/tabs", um.userManagementTabsGet).SetPermissionName(UserManagementTabsGet).Build()
	b.NewHandle().Register(builders.GET, "/users", um.userManagementUsersGet).SetPermissionName(UserManagementUsersGet).Build()
	b.NewHandle().Register(builders.GET, "/levels", um.userManagementLevelsGet).SetPermissionName(UserManagementLevelsGet).Build()
}

// func (um *userManagement) userManagementTabsGet(res http.ResponseWriter, req *http.Request) {
// 	lgr := um.Lgr("userManagementTabsGet")
// 	lgr.Info("Called")
// 	ctx := req.Context()
// 	render.Tab(req, tabModels, 0).Render(ctx, res)
// }

func (um *userManagement) userManagementUsersGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementUsersGet")
	lgr.Info("Called")
	ctx := req.Context()

	dao := um.DM().UsersDAO()
	users, err := dao.GetUserPrivilegeJoinAll()
	if err != nil || users == nil {
		tools.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}

	if len(*users) == 0 {
		datadisplay.AddTextToast(models.Warning, "No Users Found", 5).Render(ctx, res)
		return
	}

	page := pages.UserManagementUsers(*users)
	render.Tab(req, tabModels, 0, page).Render(ctx, res)
}

func (um *userManagement) userManagementLevelsGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementLevelsGet")
	lgr.Info("Called")
	ctx := req.Context()

	privileges, err := um.DM().PrivilegeDAO().GetAllJoined()
	if err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}

	page := pages.UserManagementLevels(privileges)
	render.Tab(req, tabModels, 1, page).Render(ctx, res)
}
