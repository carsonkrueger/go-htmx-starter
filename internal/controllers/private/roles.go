package private

import (
	gctx "context"
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/datadisplay"
	"github.com/carsonkrueger/main/internal/templates/datainput"
	"github.com/carsonkrueger/main/pkg/util"
	"github.com/go-chi/chi/v5"
)

const (
	RoleSelectGet  = "PrivilegeRoleSelectGet"
	SetUserRolePut = "SetUserRolePut"
)

type roles struct {
	context.AppContext
}

func NewRoles(ctx context.AppContext) *roles {
	return &roles{
		AppContext: ctx,
	}
}

func (um roles) Path() string {
	return "/roles"
}

func (um *roles) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/select", um.rolesSelectGet).SetRequiredPrivileges([]string{RoleSelectGet}).Build(ctx)
	b.NewHandler().Register(http.MethodPut, "/user/{user}", um.setUserRolePut).SetRequiredPrivileges([]string{SetUserRolePut}).Build(ctx)
}

func (r *roles) rolesSelectGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("rolesSelectGet")
	lgr.Info("Called")
	ctx := req.Context()

	defaultRole := req.URL.Query().Get("role")

	dao := r.DM().RolesDAO()
	roles, err := dao.GetAll(ctx)
	if err != nil {
		util.HandleError(req, res, lgr, err, 500, "Error fetching roles")
		return
	}

	var options []datainput.SelectOptions
	for _, r := range roles {
		options = append(options, datainput.SelectOptions{
			Value: strconv.FormatInt(int64(r.ID), 10),
			Label: r.Name,
		})
	}

	datainput.Select("roles-select", "roles", defaultRole, options, nil).Render(ctx, res)
}

func (r *roles) setUserRolePut(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("setUserRolePut")
	lgr.Info("Called")
	ctx := req.Context()

	userID := chi.URLParam(req, "user")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid user id")
		return
	}

	if err := req.ParseForm(); err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid form data")
		return
	}

	rolesIDStr := req.Form.Get("role")
	roleID, err := strconv.ParseInt(rolesIDStr, 10, 16)
	if err != nil {
		util.HandleError(req, res, lgr, err, 400, "Invalid role id")
		return
	}

	if err := r.SM().PrivilegesService().SetUserRole(ctx, int16(roleID), userIDInt); err != nil {
		util.HandleError(req, res, lgr, err, 500, "Error setting user role")
		return
	}

	if util.IsHxRequest(req) {
		datadisplay.AddTextToast(datadisplay.Success, "Success", "User Role Updated").Render(ctx, res)
	}
}
