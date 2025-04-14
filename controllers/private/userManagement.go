package private

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	tabs "github.com/carsonkrueger/main/templates/tabs/userManagement"
	"github.com/carsonkrueger/main/tools"
)

const (
	UserManagementGet       = "UserManagementGet"
	UserManagementUsersGet  = "UserManagementUsersGet"
	UserManagementLevelsGet = "UserManagementLevelsGet"
)

var tabModels = []pageLayouts.TabModel{
	{Title: "Users", PushUrl: "/user_management/users", HxGet: "/user_management/users"},
	{Title: "Privilege Levels", PushUrl: "/user_management/levels", HxGet: "/user_management/levels"},
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
	b.NewHandle().Register(builders.GET, "/tabs", um.userManagementTabsGet).SetPermissionName(UserManagementUsersGet).Build()
	b.NewHandle().Register(builders.GET, "/users", um.userManagementUsersGet).SetPermissionName(UserManagementUsersGet).Build()
	b.NewHandle().Register(builders.GET, "/levels", um.userManagementLevelsGet).SetPermissionName(UserManagementLevelsGet).Build()
}

func (um *userManagement) userManagementTabsGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementTabsGet")
	lgr.Info("userManagementTabsGet Called")
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	tabIdxStr := req.URL.Query().Get("tab")
	tabIdx, err := strconv.Atoi(tabIdxStr)
	if err != nil || tabIdx < 0 || tabIdx >= len(tabModels) {
		tabIdx = 0
	}
	um.GetTab(hxRequest, tabIdx).Render(ctx, res)
}

func (um *userManagement) userManagementUsersGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementUsersGet")
	lgr.Info("userManagementUsersGet Called")
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	if !hxRequest {
		um.GetTab(hxRequest, 0).Render(ctx, res)
		return
	}
	dao := um.DM().UsersDAO()
	users, err := dao.GetAll()
	if err != nil {
		lgr.Error(err.Error())
		datadisplay.AddToastErrors(0, err)
		return
	}
	tabs.Users(*users).Render(ctx, res)
}

func (um *userManagement) userManagementLevelsGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementLevelsGet")
	lgr.Info("userManagementLevelsGet Called")
	ctx := req.Context()
	hxRequest := tools.IsHxRequest(req)
	if !hxRequest {
		um.GetTab(hxRequest, 1).Render(ctx, res)
		return
	}
	tabs.Levels().Render(ctx, res)
}

func (um *userManagement) GetTab(hxRequest bool, index int) templ.Component {
	page := pageLayouts.MainPageLayout(pageLayouts.Tabs(tabModels, "/user_management/tabs", index))
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	return page
}
