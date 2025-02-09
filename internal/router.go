package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/private_routes"
	"github.com/carsonkrueger/main/internal/public_routes"
	"github.com/carsonkrueger/main/internal/types"
	"github.com/go-chi/chi/v5"
)

type AppRouter struct {
	public  []types.AppPublicRoute
	private []types.AppPrivateRoute
	addr    string
	router  chi.Router
}

func (a *AppRouter) Setup() {
	a.public = []types.AppPublicRoute{
		&public_routes.HelloWorld{},
		&public_routes.WebPublic{},
	}
	a.private = []types.AppPrivateRoute{
		&private_routes.HelloWorld2{},
	}
}

func (a *AppRouter) BuildRouter(ctx *types.AppContext) {
	a.router = chi.NewRouter()

	fmt.Println("Creating public routes")
	for _, r := range a.public {
		r.SetCtx(ctx)
		router := chi.NewRouter()
		r.PublicRoute(router)
		a.router.Mount(r.Path(), router)
		fmt.Printf("Registered %v\n", r.Path())
	}

	// enforce auth mw
	// router.Use()

	fmt.Println("Creating private routes")
	for _, r := range a.private {
		r.SetCtx(ctx)
		builder := builders.NewPrivateRouteBuilder()
		r.PrivateRoute(&builder)
		a.router.Mount(r.Path(), builder.Build())
		fmt.Printf("Registered %v\n", r.Path())
	}
}

func (a *AppRouter) Start(addr string, port int) error {
	if a.router == nil {
		return errors.New("AppRouter has no router. Did you forget to call BuildRouter().")
	}

	a.addr = fmt.Sprintf("%v:%v", addr, port)
	fmt.Printf("Listening on %v\n", a.addr)
	return http.ListenAndServe(a.addr, a.router)
}
