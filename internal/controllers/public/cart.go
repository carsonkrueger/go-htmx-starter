package public

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/pages"
	"github.com/carsonkrueger/main/pkg/util/render"
	"github.com/go-chi/chi/v5"
)

type cart struct {
	context.AppContext
}

func NewCart(ctx context.AppContext) *cart {
	return &cart{
		AppContext: ctx,
	}
}

func (c *cart) Path() string {
	return "/cart"
}

func (c *cart) PublicRoute(r chi.Router) {
	r.Get("/", c.cart)
}

func (c *cart) cart(res http.ResponseWriter, req *http.Request) {
	lgr := c.Lgr("cart")
	lgr.Info("Called")
	ctx := req.Context()
	page := pages.Cart()
	render.PageMainLayout(req, page).Render(ctx, res)
}
