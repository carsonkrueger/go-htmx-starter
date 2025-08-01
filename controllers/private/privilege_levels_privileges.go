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
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
)

const (
	PrivilegeLevelsPrivilegesPost   = "PrivilegeLevelsPrivilegesPost"
	PrivilegeLevelsPrivilegesDelete = "PrivilegeLevelsPrivilegesDelete"
)

type privilegeLevelsPrivileges struct {
	context.AppContext
}

func NewPrivilegeLevelsPrivileges(ctx context.AppContext) *privilegeLevelsPrivileges {
	return &privilegeLevelsPrivileges{
		AppContext: ctx,
	}
}

func (um privilegeLevelsPrivileges) Path() string {
	return "/privilege-levels-privileges"
}

func (um *privilegeLevelsPrivileges) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(builders.POST, "/", um.privilegeLevelsPrivilegesPost).SetPermissionName(PrivilegeLevelsPrivilegesPost).Build(ctx)
	b.NewHandler().Register(builders.DELETE, "/level/{level}/privilege/{privilege}", um.privilegeLevelsPrivilegesDelete).SetPermissionName(PrivilegeLevelsPrivilegesDelete).Build(ctx)
}

func (r *privilegeLevelsPrivileges) privilegeLevelsPrivilegesPost(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegeLevelsPrivilegesPost")
	lgr.Info("Called")
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid Form")
		return
	}

	level := req.FormValue("privilege-levels")
	privilege := req.FormValue("privileges")

	levelInt, err := strconv.ParseInt(level, 10, 64)
	if err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid privilege level")
		return
	}
	privilegeInt, err := strconv.ParseInt(privilege, 10, 64)
	if err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid privilege")
		return
	}

	if r.SM().PrivilegesService().HasPermissionByID(ctx, levelInt, privilegeInt) {
		tools.HandleError(req, res, lgr, err, 400, "Privilege already exists")
		return
	}

	priv, err := r.DM().PrivilegeDAO().GetOne(ctx, privilegeInt)
	if err != nil || priv == nil {
		tools.HandleError(req, res, lgr, err, 404, "Privilege not found")
		return
	}

	lvl, err := r.DM().PrivilegeLevelsDAO().GetOne(ctx, levelInt)
	if err != nil || lvl == nil {
		tools.HandleError(req, res, lgr, err, 404, "Privilege level not found")
		return
	}

	if err = r.SM().PrivilegesService().CreatePrivilegeAssociation(ctx, levelInt, priv.ID); err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Failed to add permission")
		return
	}

	if tools.IsHxRequest(req) {
		jpl := []auth_models.JoinedPrivilegesRaw{
			{
				LevelID:            levelInt,
				LevelName:          lvl.Name,
				PrivilegeID:        priv.ID,
				PrivilegeName:      priv.Name,
				PrivilegeCreatedAt: priv.CreatedAt,
			},
		}
		rows := r.SM().PrivilegesService().JoinedPrivilegesAsRowData(ctx, jpl)
		tr := datadisplay.BasicTR(rows[0])
		toast := datadisplay.AddTextToast(datadisplay.Success, "Added privilege level", 3)
		templ.Join(tr, toast).Render(ctx, res)
	}
}

func (r *privilegeLevelsPrivileges) privilegeLevelsPrivilegesDelete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegeLevelsPrivilegesDelete")
	lgr.Info("Called")
	ctx := req.Context()

	level := chi.URLParam(req, "level")
	privilege := chi.URLParam(req, "privilege")

	levelInt, err := strconv.ParseInt(level, 10, 64)
	if err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid privilege level")
		return
	}
	privilegeInt, err := strconv.ParseInt(privilege, 10, 64)
	if err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid privilege")
		return
	}

	if err := r.SM().PrivilegesService().DeletePrivilegeAssociation(ctx, levelInt, privilegeInt); err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Failed to remove permission")
		return
	}

	if tools.IsHxRequest(req) {
		datadisplay.AddTextToast(datadisplay.Success, "Deleted privilege level", 3).Render(ctx, res)
	}
}
