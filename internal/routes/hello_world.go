package routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/go-chi/chi/v5"
)

type HelloWorld struct {
	ctx *types.AppContext
}

func (r HelloWorld) Path() string {
	return "/helloworld"
}

func (r HelloWorld) SetCtx(ctx *types.AppContext) {
	r.ctx = ctx
}

func (hw HelloWorld) PublicRoute(r chi.Router) {
	// builder := builders.NewPrivateRouteBuilder()
	// builder.NewHandle().RegisterRoute("get", "/", get, enums.HelloWorldGet).Build()
	// return builder.Build()
	r.Get("/", get)
}

func get(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
