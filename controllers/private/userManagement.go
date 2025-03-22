package private

import (
	"context"
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/tools"
)

const (
	UserManagementGet       = "UserManagementGet"
	UserManagementUsersGet  = "UserManagementUsersGet"
	UserManagementLevelsGet = "UserManagementLevelsGet"
)

var tabs []pageLayouts.TabModel = []pageLayouts.TabModel{
	{Title: "Users", PushUrl: "/user_management/users", HxGet: "/user_management/users"},
	{Title: "Privilege Levels", PushUrl: "/user_management/levels", HxGet: "/user_management/levels"},
}

type UserManagement struct {
	interfaces.IAppContext
}

func (um *UserManagement) SetAppCtx(ctx interfaces.IAppContext) {
	um.IAppContext = ctx
}

func (r UserManagement) Path() string {
	return "/user_management"
}

func (um *UserManagement) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().Register(builders.GET, "/", um.userManagementGet).SetPermissionName(UserManagementGet).Build()
	b.NewHandle().Register(builders.GET, "/users", um.userManagementUsersGet).SetPermissionName(UserManagementUsersGet).Build()
	b.NewHandle().Register(builders.GET, "/levels", um.userManagementLevelsGet).SetPermissionName(UserManagementLevelsGet).Build()
	// b.NewHandle().RegisterRoute(controllers.GET, "/get2", um.hello2).SetPermission(&enums.HelloWorldGet2).Build()
}

func (um *UserManagement) userManagementGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr()
	lgr.Info("GET /user_managment")
	ctx := req.Context()
	RenderTab(res, req, ctx, 0)
}

func (um *UserManagement) userManagementUsersGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr()
	lgr.Info("GET /user_management/users")
	ctx := req.Context()
	RenderTab(res, req, ctx, 0)
}

func (um *UserManagement) userManagementLevelsGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr()
	lgr.Info("GET /user_management/levels")
	ctx := req.Context()
	RenderTab(res, req, ctx, 1)
}

func RenderTab(res http.ResponseWriter, req *http.Request, ctx context.Context, index int) error {
	hxRequest := tools.IsHxRequest(req)
	tabLayout := pageLayouts.Tabs(tabs, index)
	page := pageLayouts.MainPageLayout(tabLayout)
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	return page.Render(ctx, res)
}
