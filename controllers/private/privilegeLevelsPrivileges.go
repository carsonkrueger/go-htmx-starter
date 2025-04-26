package private

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
)

const (
	PrivilegeLevelsPrivilegesPost   = "PrivilegeLevelsPrivilegesPost"
	PrivilegeLevelsPrivilegesDelete = "PrivilegeLevelsPrivilegesDelete"
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
	b.NewHandle().Register(builders.POST, "/", um.privilegeLevelsPrivilegesPost).SetPermissionName(PrivilegeLevelsPrivilegesPost).Build()
	b.NewHandle().Register(builders.DELETE, "/level/{level}/privilege/{privilege}", um.privilegeLevelsPrivilegesDelete).SetPermissionName(PrivilegeLevelsPrivilegesDelete).Build()
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

	pk := authModels.PrivilegeLevelsPrivilegesPrimaryKey{
		PrivilegeLevelID: levelInt,
		PrivilegeID:      privilegeInt,
	}

	db := r.DB()

	if row, _ := r.DM().PrivilegeLevelsPrivilegesDAO().GetOne(pk, db); row != nil {
		tools.HandleError(req, res, lgr, errors.New("privilege level privilege already exists"), 400, "Already exists")
		return
	}

	if err := r.DM().PrivilegeLevelsPrivilegesDAO().Insert(model, db); err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Failed to insert privilege level privilege")
		return
	}

	if tools.IsHxRequest(req) {
		datadisplay.AddTextToast(models.Success, "Added privilege level", 3).Render(ctx, res)
	}
}

func (r *privilegeLevelsPrivileges) privilegeLevelsPrivilegesDelete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegeLevelsPrivilegesDelete")
	lgr.Info("Called")

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

	pk := authModels.PrivilegeLevelsPrivilegesPrimaryKey{
		PrivilegeLevelID: levelInt,
		PrivilegeID:      privilegeInt,
	}

	if err := r.DM().PrivilegeLevelsPrivilegesDAO().Delete(pk, r.DB()); err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Failed to insert privilege level privilege")
		return
	}

}
