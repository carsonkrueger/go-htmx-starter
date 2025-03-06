package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/internal"
	"github.com/carsonkrueger/main/internal/private_routes"
	"github.com/carsonkrueger/main/internal/public_routes"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

type AppRouter struct {
	public  []internal.AppPublicRoute
	private []internal.AppPrivateRoute
	addr    string
	router  chi.Router
}

func Setup() AppRouter {
	return AppRouter{
		public: []internal.AppPublicRoute{
			&public_routes.Auth{},
			&public_routes.Home{},
			&public_routes.WebPublic{},
		},
		private: []internal.AppPrivateRoute{
			&private_routes.HelloWorld2{},
		},
	}
}

func (a *AppRouter) BuildRouter(ctx *internal.AppContext) {
	a.router = chi.NewRouter()

	fmt.Println("Creating public routes:")
	for _, r := range a.public {
		r.SetAppCtx(ctx)
		router := chi.NewRouter()
		r.PublicRoute(router)
		a.router.Mount(r.Path(), router)
		fmt.Printf("\t%v\n", r.Path())
	}

	// enforce auth mw
	// router.Use()

	fmt.Println("\nCreating private routes:")
	for _, r := range a.private {
		r.SetAppCtx(ctx)
		builder := internal.NewPrivateRouteBuilder()
		r.PrivateRoute(&builder)
		a.router.Mount(r.Path(), builder.Build())
		fmt.Printf("\t%v\n", r.Path())
	}
}

func (a *AppRouter) Start(cfg cfg.Config) error {
	if a.router == nil {
		return errors.New("AppRouter has no router. Did you forget to call BuildRouter().")
	}

	a.addr = fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	fmt.Printf("\nListening on %v\n", a.addr)
	return http.ListenAndServe(a.addr, a.router)
}
