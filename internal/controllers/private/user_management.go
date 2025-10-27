package private

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/ui/layouts"
	"github.com/carsonkrueger/main/internal/templates/ui/pages"
	"github.com/carsonkrueger/main/pkg/util/render"
)

var UserManagementTabModels = []layouts.TabModel{
	{Title: "Users", PushUrl: true, HxGet: "/user_management/users"},
	{Title: "Roles", PushUrl: true, HxGet: "/user_management/roles"},
}

type userManagement struct {
	*context.AppContext
}

func NewUserManagement(ctx *context.AppContext) *userManagement {
	return &userManagement{
		AppContext: ctx,
	}
}

func (um userManagement) Path() string {
	return "/user_management"
}

func (um *userManagement) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/users", um.userManagementUsersGet).SetRequiredPrivileges(constant.UsersList).Build(ctx)
	b.NewHandler().Register(http.MethodGet, "/roles", um.userManagementRolesGet).SetRequiredPrivileges(constant.RolesList).Build(ctx)
}

func (um *userManagement) userManagementUsersGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementUsersGet")
	lgr.Info("Called")
	ctx := req.Context()

	dao := um.DM().UsersDAO()
	users, err := dao.GetUserPrivilegeJoinAll(ctx)
	if err != nil || users == nil {
		common.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}

	if len(*users) == 0 {
		return
	}

	allRoles, err := um.DM().RolesDAO().GetAll(ctx)
	if err != nil {
		common.HandleError(req, res, lgr, err, 500, "Error fetching roles")
		return
	}

	rows := um.SM().PrivilegesService().UserRoleJoinAsRowData(ctx, *users, allRoles)
	page := pages.UserManagementUsers(rows)
	render.Tab(req, UserManagementTabModels, 0, page).Render(ctx, res)
}

func (um *userManagement) userManagementRolesGet(res http.ResponseWriter, req *http.Request) {
	lgr := um.Lgr("userManagementRolesGet")
	lgr.Info("Called")
	ctx := req.Context()

	privileges, err := um.DM().PrivilegeDAO().GetAllJoined(ctx)
	if err != nil {
		common.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}
	rows := um.SM().PrivilegesService().JoinedPrivilegesAsRowData(ctx, privileges)

	page := pages.UserManagementRoles(rows)
	render.Tab(req, UserManagementTabModels, 1, page).Render(ctx, res)
}
