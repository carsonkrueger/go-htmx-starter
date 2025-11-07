package private

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/templatetargets"
	"github.com/carsonkrueger/main/internal/templates/ui/layouts"
	"github.com/carsonkrueger/main/internal/templates/ui/pages"
	"github.com/carsonkrueger/main/internal/templates/ui/tables"
	"github.com/carsonkrueger/main/pkg/util/render"
)

var UserManagementTabModels = []layouts.TabModel{
	{Title: "Users", PushUrl: true, HxGet: "/management/users"},
	{Title: "Roles", PushUrl: true, HxGet: "/management/roles"},
}

type management struct {
	*context.AppContext
}

func NewUserManagement(ctx *context.AppContext) *management {
	return &management{
		AppContext: ctx,
	}
}

func (um management) Path() string {
	return "/management"
}

func (um *management) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/users", um.managementUsersGet).SetRequiredPrivileges(constant.UsersList).Build(ctx)
	b.NewHandler().Register(http.MethodGet, "/roles", um.managementRolesGet).SetRequiredPrivileges(constant.RolesList).Build(ctx)
}

func (um *management) managementUsersGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "management.userManagementUsersGet")
	lgr.Info("Called")

	dao := um.DM().UsersDAO()
	users, err := dao.GetUserPrivilegeJoinAll(ctx)
	if err != nil || users == nil {
		common.HandleError(req, w, lgr, err, 500, "Error fetching privileges")
		return
	}

	if len(users) == 0 {
		return
	}

	roles, err := um.DM().RolesDAO().GetAll(ctx)
	if err != nil {
		common.HandleError(req, w, lgr, err, 500, "Error fetching roles")
		return
	}

	tabs := common.Tab(req, UserManagementTabModels, 0, tables.ManageUsersTable(users, roles))
	render.Layout(ctx, req, w, templatetargets.Main, tabs)
}

func (um *management) managementRolesGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "management.managementRolesGet")
	lgr.Info("Called")

	privileges, err := um.DM().PrivilegeDAO().GetAllJoined(ctx)
	if err != nil {
		common.HandleError(req, w, lgr, err, 500, "Error fetching privileges")
		return
	}

	tabs := common.Tab(req, UserManagementTabModels, 1, pages.ManagementRoles(privileges))
	render.Layout(ctx, req, w, templatetargets.Main, tabs)
}
