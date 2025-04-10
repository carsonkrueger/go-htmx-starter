package private

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
)

const (
	UserManagementGet       = "UserManagementGet"
	UserManagementUsersGet  = "UserManagementUsersGet"
	UserManagementLevelsGet = "UserManagementLevelsGet"
)

var tabs []pageLayouts.TabModel = []pageLayouts.TabModel{
	{Title: "Users", PushUrl: "/user_management/users", HxGet: "/user_management/users", Tab: pages.Signup()},
	{Title: "Privilege Levels", PushUrl: "/user_management/levels", HxGet: "/user_management/levels", Tab: pages.Signup()},
}

type userManagement struct {
	interfaces.IAppContext
}

func NewUserManagement(ctx interfaces.IAppContext) *userManagement {
	return &userManagement{
		IAppContext: ctx,
	}
}

func (r userManagement) Path() string {
	return "/user_management"
}

func (um *userManagement) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().Register(builders.GET, "/", um.userManagementGet).SetPermissionName(UserManagementGet).Build()
	b.NewHandle().Register(builders.GET, "/users", um.userManagementUsersGet).SetPermissionName(UserManagementUsersGet).Build()
	b.NewHandle().Register(builders.GET, "/levels", um.userManagementLevelsGet).SetPermissionName(UserManagementLevelsGet).Build()
	// b.NewHandle().RegisterRoute(controllers.GET, "/get2", um.hello2).SetPermission(&enums.HelloWorldGet2).Build()
}

func (um *userManagement) userManagementGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementGet")
	lgr.Info("GET /user_managment")
	ctx := req.Context()
	GetTab(res, req, ctx, 0)
}

func (um *userManagement) userManagementUsersGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementUsersGet")
	lgr.Info("GET /user_management/users")
	ctx := req.Context()
	GetTab(res, req, ctx, 0)
}

func (um *userManagement) userManagementLevelsGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementLevelsGet")
	lgr.Info("GET /user_management/levels")
	ctx := req.Context()
	GetTab(res, req, ctx, 1)
}

func GetTab(res http.ResponseWriter, req *http.Request, ctx context.Context, index int) templ.Component {
	hxRequest := tools.IsHxRequest(req)
	page := pageLayouts.MainPageLayout(pageLayouts.Tabs(tabs, index))
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	return page
}
