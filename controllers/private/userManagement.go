package private

import (
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/templates/pages"
	"github.com/carsonkrueger/main/tools"
)

const (
	UserManagementGet = "UserManagementGet"
)

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
	// b.NewHandle().RegisterRoute(controllers.GET, "/get2", um.hello2).SetPermission(&enums.HelloWorldGet2).Build()
}

func (um *UserManagement) userManagementGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr()
	lgr.Info("GET /user_managment")
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	page := pageLayouts.MainPageLayout(pages.Home())
	// If not hx request then user just arrived. Give them the index.html
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	page.Render(ctx, res)
}
