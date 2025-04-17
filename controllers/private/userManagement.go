package private

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/tools"
)

const (
	UserManagementTabsGet   = "UserManagementTabsGet"
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
	b.NewHandle().Register(builders.GET, "/tabs", um.userManagementTabsGet).SetPermissionName(UserManagementTabsGet).Build()
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
	users, err := dao.GetUserPrivilegeJoinAll()
	if err != nil {
		lgr.Error(err.Error())
		datadisplay.AddToastErrors(0, err)
		return
	}

	if users != nil {
		if len(*users) == 0 {
			datadisplay.AddTextToast(models.Warning, "No Users Found", 5).Render(ctx, res)
		}

		headers := []datadisplay.CellData{
			{
				ID:    "h-name",
				Width: 1,
				Body:  datadisplay.Text("Name", models.LG),
			},
			{
				ID:    "h-pr",
				Width: 1,
				Body:  datadisplay.Text("Privilege Level", models.LG),
			},
			{
				ID:    "h-ca",
				Width: 1,
				Body:  datadisplay.Text("Created At", models.LG),
			},
		}

		cells := make([][]datadisplay.CellData, len(*users))
		for i, user := range *users {
			cells[i] = []datadisplay.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(user.FirstName, models.MD),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(strconv.FormatInt(user.PrivilegeLevelID, 10), models.MD),
				},
				{
					ID:    "ca-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(user.CreatedAt.Format("2006-01-02"), models.MD),
				},
			}
		}
		datadisplay.BasicTable(headers, cells).Render(ctx, res)
	}
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

	// cells := make([][]datadisplay.CellData, len(*users))
	// 		for i, user := range *users {
	// 			cells[i] = []datadisplay.CellData{
	// 				{
	// 					ID:    "n-" + strconv.Itoa(i),
	// 					Width: 1,
	// 					Body:  datadisplay.Text(user.FirstName, models.MD),
	// 				},
	// 				{
	// 					ID:    "pr-" + strconv.Itoa(i),
	// 					Width: 1,
	// 					Body:  datadisplay.Text(strconv.FormatInt(user.PrivilegeLevelID, 10), models.MD),
	// 				},
	// 				{
	// 					ID:    "ca-" + strconv.Itoa(i),
	// 					Width: 1,
	// 					Body:  datadisplay.Text(user.CreatedAt.Format("2006-01-02"), models.MD),
	// 				},
	// 			}
	// 		}

	headers := []datadisplay.CellData{}
	cells := [][]datadisplay.CellData{}
	datadisplay.BasicTable(headers, cells).Render(ctx, res)
}

func (um *userManagement) GetTab(hxRequest bool, index int) templ.Component {
	page := pageLayouts.MainPageLayout(pageLayouts.Tabs(tabModels, "/user_management/tabs", index))
	if !hxRequest {
		page = pageLayouts.Index(page)
	}
	return page
}
