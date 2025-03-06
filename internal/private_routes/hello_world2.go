package private_routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal"
	"github.com/carsonkrueger/main/internal/enums"
)

type HelloWorld2 struct {
	internal.WithAppContext
}

func (r HelloWorld2) Path() string {
	return "/helloworld"
}

func (hw *HelloWorld2) PrivateRoute(b *internal.PrivateRouteBuilder) {
	b.NewHandle().RegisterRoute("get", "/", hw.hello).SetPermission(enums.HelloWorldGet).Build()
}

func (hw *HelloWorld2) hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
