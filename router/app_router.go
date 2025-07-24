package router

import (
	gctx "context"
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/controllers/private"
	"github.com/carsonkrueger/main/controllers/public"
	"github.com/carsonkrueger/main/middlewares"

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

func NewAppRouter(appCtx context.AppContext) AppRouter {
	return AppRouter{
		appCtx: appCtx,
		public: []builders.AppPublicRoute{
			// DB-START
			public.NewLogin(appCtx),
			public.NewSignUp(appCtx),
			// DB-END
			public.NewWebPublic(appCtx),
			public.NewHome(appCtx),
			// INSERT PUBLIC
		},
		// DB-START
		private: []builders.AppPrivateRoute{
			private.NewUserManagement(appCtx),
			private.NewPrivileges(appCtx),
			private.NewPrivilegeLevels(appCtx),
			private.NewPrivilegeLevelsPrivileges(appCtx),
			// INSERT PRIVATE
		},
		// DB-END
	}
}

func (a *AppRouter) BuildRouter(ctx gctx.Context) {
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
		r.PrivateRoute(ctx, &builder)
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
