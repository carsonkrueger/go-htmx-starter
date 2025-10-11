package private

import (
	gctx "context"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/util"
	"github.com/go-chi/chi/v5"
)

const (
	RolesPrivilegesPost   = "RolesPrivilegesPost"
	RolesPrivilegesDelete = "RolesPrivilegesDelete"
)

type rolesPrivileges struct {
	context.AppContext
}

func NewRolesPrivileges(ctx context.AppContext) *rolesPrivileges {
	return &rolesPrivileges{
		AppContext: ctx,
	}
}

func (um rolesPrivileges) Path() string {
	return "/roles-privileges"
}

func (um *rolesPrivileges) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(builders.POST, "/", um.rolesPrivilegesPost).SetRequiredPrivileges([]string{RolesPrivilegesPost}).Build(ctx)
	b.NewHandler().Register(builders.DELETE, "/role/{role}/privilege/{privilege}", um.rolesPrivilegesDelete).SetRequiredPrivileges([]string{RolesPrivilegesDelete}).Build(ctx)
}

func (r *rolesPrivileges) rolesPrivilegesPost(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("rolesPrivilegesPost")
	lgr.Info("Called")
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid Form")
		return
	}

	roles := req.FormValue("roles")
	privilege := req.FormValue("privileges")

	roleID, err := strconv.ParseInt(roles, 10, 64)
	if err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid roles")
		return
	}
	privilegeInt, err := strconv.ParseInt(privilege, 10, 64)
	if err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid privilege")
		return
	}

	if r.SM().PrivilegesService().HasPermissionsByIDS(ctx, int16(roleID), []int64{privilegeInt}) {
		util.HandleError(req, res, lgr, err, 400, "Privilege already exists")
		return
	}

	priv, err := r.DM().PrivilegeDAO().GetOne(ctx, privilegeInt)
	if err != nil {
		util.HandleError(req, res, lgr, err, 404, "Privilege not found")
		return
	}

	role, err := r.DM().RolesDAO().GetOne(ctx, int16(roleID))
	if err != nil {
		util.HandleError(req, res, lgr, err, 404, "Roles not found")
		return
	}

	if err = r.SM().PrivilegesService().CreatePrivilegeAssociation(ctx, int16(roleID), priv.ID); err != nil {
		util.HandleError(req, res, lgr, err, 500, "Failed to add permission")
		return
	}

	if util.IsHxRequest(req) {
		jpl := []auth_models.JoinedPrivilegesRaw{
			{
				RoleID:             int16(roleID),
				RoleName:           role.Name,
				PrivilegeID:        priv.ID,
				PrivilegeName:      priv.Name,
				PrivilegeCreatedAt: priv.CreatedAt,
			},
		}
		rows := r.SM().PrivilegesService().JoinedPrivilegesAsRowData(ctx, jpl)
		tr := datadisplay.BasicTR(rows[0])
		toast := datadisplay.AddTextToast(datadisplay.Success, "Success", "Added roles")
		templ.Join(tr, toast).Render(ctx, res)
	}
}

func (r *rolesPrivileges) rolesPrivilegesDelete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("rolesPrivilegesDelete")
	lgr.Info("Called")
	ctx := req.Context()

	role := chi.URLParam(req, "role")
	privilege := chi.URLParam(req, "privilege")

	roleID, err := strconv.ParseInt(role, 10, 64)
	if err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid role")
		return
	}
	privilegeInt, err := strconv.ParseInt(privilege, 10, 64)
	if err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid privilege")
		return
	}

	if err := r.SM().PrivilegesService().DeletePrivilegeAssociation(ctx, int16(roleID), privilegeInt); err != nil {
		util.HandleError(req, res, lgr, err, 500, "Failed to remove permission")
		return
	}

	if util.IsHxRequest(req) {
		datadisplay.AddTextToast(datadisplay.Success, "Success", "Deleted role").Render(ctx, res)
	}
}
