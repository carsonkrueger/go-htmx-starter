package private

import (
	gctx "context"
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/toast"
	"github.com/carsonkrueger/main/internal/templates/ui/tables"
	"github.com/carsonkrueger/main/pkg/model"
	tuitoast "github.com/carsonkrueger/main/pkg/templui/toast"
	"github.com/carsonkrueger/main/pkg/util"
	"github.com/go-chi/chi/v5"
)

type rolesPrivileges struct {
	*context.AppContext
}

func NewRolesPrivileges(ctx *context.AppContext) *rolesPrivileges {
	return &rolesPrivileges{
		AppContext: ctx,
	}
}

func (um rolesPrivileges) Path() string {
	return "/roles-privileges"
}

func (um *rolesPrivileges) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodPost, "/", um.rolesPrivilegesPost).SetRequiredPrivileges(constant.RolesUpdate, constant.PrivilegesUpdate).Build(ctx)
	b.NewHandler().Register(http.MethodDelete, "/role/{role}/privilege/{privilege}", um.rolesPrivilegesDelete).SetRequiredPrivileges(constant.RolesUpdate, constant.PrivilegesUpdate).Build(ctx)
}

func (r *rolesPrivileges) rolesPrivilegesPost(w http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("rolesPrivilegesPost")
	lgr.Info("Called")
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid Form")
		return
	}

	roles := req.FormValue("roles")
	privilege := req.FormValue("privileges")

	roleID, err := strconv.ParseInt(roles, 10, 64)
	if err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid roles")
		return
	}
	privilegeInt, err := strconv.ParseInt(privilege, 10, 64)
	if err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid privilege")
		return
	}

	if r.SM().PrivilegesService().HasPermissionsByIDS(ctx, int16(roleID), []int64{privilegeInt}) {
		common.HandleError(req, w, lgr, err, 400, "Privilege already exists")
		return
	}

	priv, err := r.DM().PrivilegeDAO().GetOne(ctx, privilegeInt)
	if err != nil {
		common.HandleError(req, w, lgr, err, 404, "Privilege not found")
		return
	}

	role, err := r.DM().RolesDAO().GetOne(ctx, int16(roleID))
	if err != nil {
		common.HandleError(req, w, lgr, err, 404, "Roles not found")
		return
	}

	if err = r.SM().PrivilegesService().CreatePrivilegeAssociation(ctx, int16(roleID), priv.ID); err != nil {
		common.HandleError(req, w, lgr, err, 500, "Failed to add permission")
		return
	}

	if util.IsHxRequest(req) {
		tables.RolesRow(model.RolesPrivilegeJoin{Privileges: priv, Roles: role}).Render(ctx, w)
		toast.ToastsSwap(tuitoast.Props{Title: "Success", Description: "Privilege added to role"}).Render(ctx, w)
	}
}

func (r *rolesPrivileges) rolesPrivilegesDelete(w http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("rolesPrivilegesDelete")
	lgr.Info("Called")
	ctx := req.Context()

	role := chi.URLParam(req, "role")
	privilege := chi.URLParam(req, "privilege")

	roleID, err := strconv.ParseInt(role, 10, 64)
	if err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid role")
		return
	}
	privilegeInt, err := strconv.ParseInt(privilege, 10, 64)
	if err != nil {
		common.HandleError(req, w, lgr, err, 400, "Invalid privilege")
		return
	}

	if err := r.SM().PrivilegesService().DeletePrivilegeAssociation(ctx, int16(roleID), privilegeInt); err != nil {
		common.HandleError(req, w, lgr, err, 500, "Failed to remove permission")
		return
	}

	if util.IsHxRequest(req) {
		toast.ToastsSwap(tuitoast.Props{Title: "Success", Description: "Privilege removed from role"}).Render(ctx, w)
	}
}
