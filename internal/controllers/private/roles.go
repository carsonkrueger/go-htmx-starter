package private

import (
	gctx "context"
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/selectbox"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/toast"
	tuitoast "github.com/carsonkrueger/main/pkg/templui/toast"
	"github.com/carsonkrueger/main/pkg/util"
	"github.com/go-chi/chi/v5"
)

type roles struct {
	*context.AppContext
}

func NewRoles(ctx *context.AppContext) *roles {
	return &roles{
		AppContext: ctx,
	}
}

func (um roles) Path() string {
	return "/roles"
}

func (um *roles) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/select", um.rolesSelectGet).SetRequiredPrivileges(constant.RolesList).Build(ctx)
	b.NewHandler().Register(http.MethodPut, "/user/{user}", um.setUserRolePut).SetRequiredPrivileges(constant.UsersUpdate, constant.RolesUpdate).Build(ctx)
}

func (r *roles) rolesSelectGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "roles.rolesSelectGet")
	lgr.Info("Called")

	selected := req.URL.Query().Get("role")
	selectedID, _ := strconv.ParseInt(selected, 10, 64)

	dao := r.DM().RolesDAO()
	roles, err := dao.GetAll(ctx)
	if err != nil {
		common.HandleError(req, w, lgr, err, 500, "Error fetching roles")
		return
	}

	selectbox.RoleSelectBox(selectbox.RoleSelectBoxProps{Roles: roles, SelectedRoleID: int16(selectedID)}).Render(ctx, w)
}

func (r *roles) setUserRolePut(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "roles.setUserRolePut")
	lgr.Info("Called")

	userID := chi.URLParam(req, "user")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid user id")
		return
	}

	if err := req.ParseForm(); err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid form data")
		return
	}

	rolesIDStr := req.Form.Get("role")
	roleID, err := strconv.ParseInt(rolesIDStr, 10, 16)
	if err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid role id")
		return
	}

	if err := r.SM().PrivilegesService().SetUserRole(ctx, int16(roleID), userIDInt); err != nil {
		common.HandleError(req, w, lgr, err, 500, "Error setting user role")
		return
	}

	if util.IsHxRequest(req) {
		toast.ToastsSwap(tuitoast.Props{Title: "Success", Description: "Role assigned to user"}).Render(ctx, w)
	}
}
