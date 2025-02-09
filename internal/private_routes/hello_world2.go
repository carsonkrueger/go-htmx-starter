package private_routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/enums"
	"github.com/carsonkrueger/main/internal/types"
)

type HelloWorld2 struct {
	types.WithContext
}

func (r HelloWorld2) Path() string {
	return "/helloworld2"
}

func (hw *HelloWorld2) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().RegisterRoute("get", "/", hw.get2, enums.HelloWorldGet).Build()
	b.NewHandle().RegisterRoute("get", "/test", hw.get3, enums.HelloWorldGet).Build()
}

func (hw *HelloWorld2) get2(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World2!"))
}

func (hw *HelloWorld2) get3(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World3!"))
}
