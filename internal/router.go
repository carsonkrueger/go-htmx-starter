package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/internal/routes"
	"github.com/carsonkrueger/main/internal/types"
	"github.com/go-chi/chi/v5"
)

type AppRouter struct {
	public  []types.AppRoute
	private []types.AppRoute
	addr    string
	router  chi.Router
}

func (a *AppRouter) Setup(ctx *types.AppContext) {
	a.public = []types.AppRoute{
		routes.HelloWorld{Ctx: ctx},
	}
	a.private = []types.AppRoute{}
}

func (a *AppRouter) BuildRouter() {
	a.router = chi.NewRouter()

	fmt.Println("Creating public routes")
	for _, r := range a.public {
		fmt.Printf("Registered %v\n", r.Path())
		a.router.Mount(r.Path(), r.Route())
	}

	// enforce auth mw
	// router.Use()

	fmt.Println("Creating private routes")
	for _, r := range a.private {
		fmt.Printf("Registered %v\n", r.Path())
		a.router.Mount(r.Path(), r.Route())
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
