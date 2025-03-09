package private

import (
	"net/http"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/enums"
	"github.com/carsonkrueger/main/routes"
)

type HelloWorld struct {
	context.WithAppContext
}

func (r HelloWorld) Path() string {
	return "/helloworld"
}

func (hw *HelloWorld) PrivateRoute(b *routes.PrivateRouteBuilder) {
	b.NewHandle().RegisterRoute("get", "/", hw.hello).SetPermission(&enums.HelloWorldGet).Build()
	b.NewHandle().RegisterRoute("get", "/get2", hw.hello2).SetPermission(&enums.HelloWorldGet2).Build()
}

func (hw *HelloWorld) hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}

func (hw *HelloWorld) hello2(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Secret Hello World!"))
}
