package router

import (
	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/controllers/private"
	"github.com/carsonkrueger/main/controllers/public"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/middlewares"
	"go.uber.org/zap"

	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

type AppRouter struct {
	public  []builders.IAppPublicRoute
	private []builders.IAppPrivateRoute
	addr    string
	router  chi.Router
}

func Setup(ctx interfaces.IAppContext) AppRouter {
	return AppRouter{
		public: []builders.IAppPublicRoute{
			public.NewWebPublic(ctx),
			public.NewLogin(ctx),
			public.NewSignUp(ctx),
			public.NewHome(ctx),
		},
		private: []builders.IAppPrivateRoute{
			private.NewUserManagement(ctx),
		},
	}
}

func (a *AppRouter) BuildRouter(ctx interfaces.IAppContext) {
	a.router = chi.NewRouter()
	lgr := ctx.Lgr("BuildRouter")

	for _, r := range a.public {
		router := chi.NewRouter()
		r.PublicRoute(router)
		a.router.Mount(r.Path(), router)
		lgr.Info(r.Path())
	}

	// enforce authentication middleware
	a.router = a.router.With(middlewares.EnforceAuth(ctx))

	for _, r := range a.private {
		builder := builders.NewPrivateRouteBuilder(ctx)
		r.PrivateRoute(&builder)
		a.router.Mount(r.Path(), builder.Build())
		lgr.Info(r.Path())
	}

	err := ctx.SM().PrivilegesService().BuildCache()
	if err != nil {
		lgr.Fatal("Error building permission cache", zap.Error(err))
	}
}

func (a *AppRouter) Start(cfg cfg.Config) error {
	if a.router == nil {
		return errors.New("AppRouter has no router. Did you forget to call BuildRouter().")
	}

	a.addr = fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	fmt.Printf("\nListening on http://%v\n", a.addr)
	return http.ListenAndServe(a.addr, a.router)
}
