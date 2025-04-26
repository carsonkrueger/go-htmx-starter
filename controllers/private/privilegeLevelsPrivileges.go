package private

import (
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/tools"
)

const (
	// PrivilegeLevelsPrivilegesGet   = "PrivilegeLevelsPrivilegesGet"
	PrivilegeLevelsPrivilegesPost = "PrivilegeLevelsPrivilegesPost"
	// PrivilegeLevelsPrivilegesPut = "PrivilegeLevelsPrivilegesPut"
	// PrivilegeLevelsPrivilegesPatch = "PrivilegeLevelsPrivilegesPatch"
	// PrivilegeLevelsPrivilegesDelete = "PrivilegeLevelsPrivilegesDelete"
)

type privilegeLevelsPrivileges struct {
	interfaces.IAppContext
}

func NewPrivilegeLevelsPrivileges(ctx interfaces.IAppContext) *privilegeLevelsPrivileges {
	return &privilegeLevelsPrivileges{
		IAppContext: ctx,
	}
}

func (um privilegeLevelsPrivileges) Path() string {
	return "/privilege-levels-privileges"
}

func (um *privilegeLevelsPrivileges) PrivateRoute(b *builders.PrivateRouteBuilder) {
	// b.NewHandle().Register(builders.GET, "/", um.privilegeLevelsPrivilegesGet).SetPermissionName(PrivilegeLevelsPrivilegesGet).Build()
	b.NewHandle().Register(builders.POST, "/", um.privilegeLevelsPrivilegesPost).SetPermissionName(PrivilegeLevelsPrivilegesPost).Build()
	// b.NewHandle().Register(builders.PUT, "/", um.privilegeLevelsPrivilegesPut).SetPermissionName(PrivilegeLevelsPrivilegesPut).Build()
	// b.NewHandle().Register(builders.PATCH, "/", um.privilegeLevelsPrivilegesPatch).SetPermissionName(PrivilegeLevelsPrivilegesPatch).Build()
	// b.NewHandle().Register(builders.DELETE, "/", um.privilegeLevelsPrivilegesDelete).SetPermissionName(PrivilegeLevelsPrivilegesDelete).Build()
}

func (r *privilegeLevelsPrivileges) privilegeLevelsPrivilegesGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegeLevelsPrivilegesGet")
	lgr.Info("Called")
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

	model := &model.PrivilegeLevelsPrivileges{
		PrivilegeLevelID: levelInt,
		PrivilegeID:      privilegeInt,
	}

	if err := r.DM().PrivilegeLevelsPrivilegesDAO().Insert(model, r.DB()); err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Failed to insert privilege level privilege")
		return
	}

	if tools.IsHxRequest(req) {
		datadisplay.AddTextToast(models.Success, "Added privilege level", 3).Render(ctx, res)
	}
}

func (r *privilegeLevelsPrivileges) privilegeLevelsPrivilegesPut(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegeLevelsPrivilegesPut")
	lgr.Info("Called")
}

func (r *privilegeLevelsPrivileges) privilegeLevelsPrivilegesPatch(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegeLevelsPrivilegesPatch")
	lgr.Info("Called")
}

func (r *privilegeLevelsPrivileges) privilegeLevelsPrivilegesDelete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegeLevelsPrivilegesDelete")
	lgr.Info("Called")
}
