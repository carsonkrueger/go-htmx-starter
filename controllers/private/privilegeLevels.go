package private

import (
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/templates/datainput"
	"github.com/carsonkrueger/main/tools"
)

const (
	PrivilegeLevelsSelectGet = "PrivilegeLevelsSelectGet"
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
}

func (r *privilegeLevels) privilegeLevelsSelectGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegesLevelsSelectGet")
	lgr.Info("Called")
	ctx := req.Context()

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

	datainput.Select("privileges-levels-select", "privilege-levels", options).Render(ctx, res)
}
