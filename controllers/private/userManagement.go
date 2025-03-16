package private

import (
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
)

type UserManagement struct {
	interfaces.IAppContext
}

func (um *UserManagement) SetAppCtx(ctx interfaces.IAppContext) {
	um.IAppContext = ctx
}

func (r UserManagement) Path() string {
	return "/user_management"
}

func (um *UserManagement) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().Register(builders.GET, "/", um.hello).SetPermissionName("HelloWorldGet").Build()
	// b.NewHandle().RegisterRoute(controllers.GET, "/get2", um.hello2).SetPermission(&enums.HelloWorldGet2).Build()
}

func (um *UserManagement) hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
