package main

import (
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/internal"
	"github.com/carsonkrueger/main/internal/logger"
	"github.com/carsonkrueger/main/internal/services"

	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/internal/private_routes"
	"github.com/carsonkrueger/main/internal/public_routes"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

func main() {
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)

	db, err := sql.Open("postgres", cfg.DbUrl())
	// defer db.Close()
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("Database connection is nil")
	}
	sm := services.NewServiceManager(db)

	ctx := internal.NewAppContext(lgr, sm)
	// defer ctx.CleanUp()

	appRouter := Setup()
	appRouter.BuildRouter(ctx)
	err = appRouter.Start(cfg)

	if err != nil {
		panic(err)
	}
}

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
			&private_routes.HelloWorld{},
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

	// enforce authentication middleware
	a.router = a.router.With(internal.EnforceAuth(ctx))

	fmt.Println("\nCreating private routes:")
	for _, r := range a.private {
		r.SetAppCtx(ctx)
		builder := internal.NewPrivateRouteBuilder(ctx)
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
