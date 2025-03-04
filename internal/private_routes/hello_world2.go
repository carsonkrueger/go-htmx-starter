package private_routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/enums"
	"github.com/carsonkrueger/main/internal/types"
)

type HelloWorld2 struct {
	types.WithAppContext
}

func (r HelloWorld2) Path() string {
	return "/helloworld"
}

func (hw *HelloWorld2) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().RegisterRoute("get", "/", hw.hello).SetPermission(enums.HelloWorldGet).Build()
}

func (hw *HelloWorld2) hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
