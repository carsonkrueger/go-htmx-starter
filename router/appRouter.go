package router

import (
	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/context"
	routes "github.com/carsonkrueger/main/controllers"
	"github.com/carsonkrueger/main/controllers/private"
	"github.com/carsonkrueger/main/controllers/public"
	"github.com/carsonkrueger/main/middlewares"

	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

type AppRouter struct {
	public  []routes.AppPublicRoute
	private []routes.AppPrivateRoute
	addr    string
	router  chi.Router
}

func Setup() AppRouter {
	return AppRouter{
		public: []routes.AppPublicRoute{
			&public.Auth{},
			&public.Home{},
			&public.WebPublic{},
		},
		private: []routes.AppPrivateRoute{
			&private.HelloWorld{},
		},
	}
}

func (a *AppRouter) BuildRouter(ctx *context.AppContext) {
	a.router = chi.NewRouter()

	fmt.Println("Creating public routes:")
	for _, r := range a.public {
		r.SetAppCtx(ctx)
		router := chi.NewRouter()
		r.PublicRoute(router)
		a.router.Mount(r.Path(), router)
		fmt.Printf("\t%v\n", r.Path())
	}

	// enforce authentication middleware
	a.router = a.router.With(middlewares.EnforceAuth(ctx))

	fmt.Println("\nCreating private routes:")
	for _, r := range a.private {
		r.SetAppCtx(ctx)
		builder := routes.NewPrivateRouteBuilder(ctx)
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
