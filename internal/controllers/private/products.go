package private

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/context"
)

const (
	ProductsGet    = "products:get"
	ProductsCreate = "products:create"
	ProductsUpdate = "products:update"
	ProductsDelete = "products:delete"
)

type products struct {
	context.AppContext
}

func NewProducts(ctx context.AppContext) *products {
	return &products{
		AppContext: ctx,
	}
}

func (r products) Path() string {
	return "/products"
}

func (r *products) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/", r.productsGet).SetRequiredPrivileges([]string{ProductsGet}).Build(ctx)
	b.NewHandler().Register(http.MethodPost, "/", r.productsPost).SetRequiredPrivileges([]string{ProductsCreate}).Build(ctx)
	b.NewHandler().Register(http.MethodPut, "/", r.productsPut).SetRequiredPrivileges([]string{ProductsUpdate}).Build(ctx)
	b.NewHandler().Register(http.MethodPatch, "/", r.productsPatch).SetRequiredPrivileges([]string{ProductsUpdate}).Build(ctx)
	b.NewHandler().Register(http.MethodDelete, "/", r.productsDelete).SetRequiredPrivileges([]string{ProductsDelete}).Build(ctx)
}

func (r *products) productsGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("productsGet")
	lgr.Info("Called")
}

func (r *products) productsPost(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("productsPost")
	lgr.Info("Called")
}

func (r *products) productsPut(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("productsPut")
	lgr.Info("Called")
}

func (r *products) productsPatch(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("productsPatch")
	lgr.Info("Called")
}

func (r *products) productsDelete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("productsDelete")
	lgr.Info("Called")
}
