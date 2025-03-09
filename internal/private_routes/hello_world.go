package private_routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal"
)

type HelloWorld struct {
	internal.WithAppContext
}

func (r HelloWorld) Path() string {
	return "/helloworld"
}

func (hw *HelloWorld) PrivateRoute(b *internal.PrivateRouteBuilder) {
	b.NewHandle().RegisterRoute("get", "/", hw.hello).SetPermission(&internal.HelloWorldGet).Build()
	b.NewHandle().RegisterRoute("get", "/get2", hw.hello2).SetPermission(&internal.HelloWorldGet2).Build()
}

func (hw *HelloWorld) hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}

func (hw *HelloWorld) hello2(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Secret Hello World!"))
}
