package private

import (
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/datainput"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
)

const (
	PrivilegeLevelsSelectGet = "PrivilegeLevelsSelectGet"
	SetUserLevelPut          = "SetUserLevelPut"
)

type privilegeLevels struct {
	context.AppContext
}

func NewPrivilegeLevels(ctx context.AppContext) *privilegeLevels {
	return &privilegeLevels{
		AppContext: ctx,
	}
}

func (um privilegeLevels) Path() string {
	return "/privilege-levels"
}

func (um *privilegeLevels) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().Register(builders.GET, "/select", um.privilegeLevelsSelectGet).SetPermissionName(PrivilegeLevelsSelectGet).Build()
	b.NewHandle().Register(builders.PUT, "/user/{user}", um.setUserLevelPut).SetPermissionName(SetUserLevelPut).Build()
}

func (r *privilegeLevels) privilegeLevelsSelectGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegesLevelsSelectGet")
	lgr.Info("Called")
	ctx := req.Context()

	defaultLevel := req.URL.Query().Get("level")

	dao := r.DM().PrivilegeLevelsDAO()
	levels, err := dao.Index(nil, r.DB())
	if err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Error fetching privilege levels")
		return
	}

	var options []datainput.SelectOptions
	if levels != nil {
		for _, lvl := range levels {
			options = append(options, datainput.SelectOptions{
				Value: strconv.FormatInt(lvl.ID, 10),
				Label: lvl.Name,
			})
		}
	}

	datainput.Select("privileges-levels-select", "privilege-levels", defaultLevel, options, nil).Render(ctx, res)
}

func (r *privilegeLevels) setUserLevelPut(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("setUserLevelPut")
	lgr.Info("Called")
	ctx := req.Context()

	userID := chi.URLParam(req, "user")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid user id")
		return
	}

	if err := req.ParseForm(); err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid form data")
		return
	}

	levelIDStr := req.Form.Get("privilege-level")
	levelID, err := strconv.ParseInt(levelIDStr, 10, 64)
	if err != nil {
		tools.HandleError(req, res, lgr, err, 400, "Invalid privilege level id")
		return
	}

	if err := r.SM().PrivilegesService().SetUserPrivilegeLevel(levelID, userIDInt); err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Error setting user privilege level")
		return
	}

	if tools.IsHxRequest(req) {
		datadisplay.AddTextToast(datadisplay.Success, "User Level Updated", 3).Render(ctx, res)
	}
}
