package private

import (
	gctx "context"
	"net/http"
	"strconv"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/common"
	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/selectbox"
)

type privileges struct {
	*context.AppContext
}

func NewPrivileges(ctx *context.AppContext) *privileges {
	return &privileges{
		AppContext: ctx,
	}
}

func (um privileges) Path() string {
	return "/privileges"
}

func (um *privileges) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/select", um.privilegesSelectGet).SetRequiredPrivileges(constant.PrivilegesList).Build(ctx)
}

func (r *privileges) privilegesSelectGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("privilegesSelectGet")
	lgr.Info("Called")
	ctx := req.Context()

	dao := r.DM().PrivilegeDAO()
	privileges, err := dao.GetAll(ctx)
	if err != nil {
		common.HandleError(req, res, lgr, err, 500, "Error fetching privileges")
		return
	}

	var options []selectbox.SelectOptions
	for _, p := range privileges {
		options = append(options, selectbox.SelectOptions{
			Value: strconv.FormatInt(p.ID, 10),
			Label: p.Name,
		})
	}

	selectbox.Select("privileges-select", "privileges", "", options, nil).Render(ctx, res)
}
