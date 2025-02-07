package routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/types"
	"github.com/go-chi/chi/v5"
)

type HelloWorld struct {
	Ctx *types.AppContext
}

func (r HelloWorld) Path() string {
	return "/helloworld"
}

func (r HelloWorld) Route() chi.Router {
	builder := builders.NewPrivateRouteBuilder("/helloworld")
	builder.NewHandle().RegisterRoute("get", "/", get, types.HelloWorldGet).Build()
	return builder.Build()
}

func get(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
