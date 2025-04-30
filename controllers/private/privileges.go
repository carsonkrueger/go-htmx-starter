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
	PrivilegesSelectGet  = "PrivilegesSelectGet"
	PrivilegesSelectPost = "PrivilegesSelectPost"
)

type privileges struct {
	context.AppContext
}

func NewPrivileges(ctx context.AppContext) *privileges {
	return &privileges{
		AppContext: ctx,
	}
}

func (um privileges) Path() string {
	return "/privileges"
}

func (um *privileges) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().Register(builders.GET, "/select", um.privilegesSelectGet).SetPermissionName(PrivilegesSelectGet).Build()
}

func (r *privileges) privilegesSelectGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegesSelectGet")
	lgr.Info("Called")
	ctx := req.Context()

	dao := r.DM().PrivilegeDAO()
	levels, err := dao.Index(nil, r.DB())
	if err != nil {
		tools.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
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

	datainput.Select("privileges-select", "privileges", "", options, nil).Render(ctx, res)
}
