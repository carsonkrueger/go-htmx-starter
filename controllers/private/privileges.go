package private

import (
	gctx "context"
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/templates/datainput"
	"github.com/carsonkrueger/main/util"
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

func (um *privileges) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(builders.GET, "/select", um.privilegesSelectGet).SetRequiredPrivileges([]string{PrivilegesSelectGet}).Build(ctx)
}

func (r *privileges) privilegesSelectGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegesSelectGet")
	lgr.Info("Called")
	ctx := req.Context()

	dao := r.DM().PrivilegeDAO()
	privileges, err := dao.GetAll(ctx)
	if err != nil {
		util.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}

	var options []datainput.SelectOptions
	for _, p := range privileges {
		options = append(options, datainput.SelectOptions{
			Value: strconv.FormatInt(p.ID, 10),
			Label: p.Name,
		})
	}

	datainput.Select("privileges-select", "privileges", "", options, nil).Render(ctx, res)
}
