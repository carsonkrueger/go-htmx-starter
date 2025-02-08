package routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/enums"
	"github.com/carsonkrueger/main/internal/types"
)

type HelloWorld2 struct {
	ctx *types.AppContext
}

func (r HelloWorld2) Path() string {
	return "/helloworld2"
}

func (r HelloWorld2) SetCtx(ctx *types.AppContext) {
	r.ctx = ctx
}

func (hw HelloWorld2) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().RegisterRoute("get", "/", get2, enums.HelloWorldGet).Build()
	b.NewHandle().RegisterRoute("get", "/test", get3, enums.HelloWorldGet).Build()
}

func get2(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World2!"))
}

func get3(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World3!"))
}
