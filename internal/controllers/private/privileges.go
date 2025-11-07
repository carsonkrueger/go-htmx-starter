package private

import (
	gctx "context"
	"net/http"

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

func (r *privileges) privilegesSelectGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "privileges.privilegesSelectGet")
	lgr.Info("Called")

	dao := r.DM().PrivilegeDAO()
	privileges, err := dao.GetAll(ctx)
	if err != nil {
		common.HandleError(req, w, lgr, err, 500, "Error fetching privileges")
		return
	}

	selectbox.PrivilegeSelectBox(selectbox.PrivilegeSelectBoxProps{Privileges: privileges}).Render(ctx, w)
}
