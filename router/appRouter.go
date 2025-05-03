package router

import (
	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/context"
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
	public []builders.AppPublicRoute
	// DB-START
	private []builders.AppPrivateRoute
	// DB-END
	addr   string
	router chi.Router
	appCtx context.AppContext
}

func NewAppRouter(ctx context.AppContext) AppRouter {
	return AppRouter{
		appCtx: ctx,
		public: []builders.AppPublicRoute{
			// DB-START
			public.NewLogin(ctx),
			public.NewSignUp(ctx),
			// DB-END
			public.NewWebPublic(ctx),
			public.NewHome(ctx),
		},
		// DB-START
		private: []builders.AppPrivateRoute{
			private.NewUserManagement(ctx),
			private.NewPrivileges(ctx),
			private.NewPrivilegeLevels(ctx),
			private.NewPrivilegeLevelsPrivileges(ctx),
		},
		// DB-END
	}
}

func (a *AppRouter) BuildRouter() {
	a.router = chi.NewRouter()
	lgr := a.appCtx.Lgr("BuildRouter")

	for _, r := range a.public {
		router := chi.NewRouter()
		r.PublicRoute(router)
		a.router.Mount(r.Path(), router)
		lgr.Info(r.Path())
	}

	// DB-START
	// enforce authentication middleware
	a.router = a.router.With(middlewares.EnforceAuth(a.appCtx))

	for _, r := range a.private {
		builder := builders.NewPrivateRouteBuilder(a.appCtx)
		r.PrivateRoute(&builder)
		a.router.Mount(r.Path(), builder.RawRouter())
		lgr.Info(r.Path())
	}
	// DB-END
}

func (a *AppRouter) Start(cfg cfg.Config) error {
	if a.router == nil {
		return errors.New("AppRouter has no router. Did you forget to call BuildRouter().")
	}

	a.addr = fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	fmt.Printf("\nListening on http://%s\n", a.addr)
	return http.ListenAndServe(a.addr, a.router)
}
