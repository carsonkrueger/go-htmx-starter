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
// DB-START
	private []builders.IAppPrivateRoute
// DB-END
	addr    string
	router  chi.Router
	appCtx  interfaces.IAppContext
}

func NewAppRouter(ctx interfaces.IAppContext) AppRouter {
	return AppRouter{
		appCtx: ctx,
		public: []builders.IAppPublicRoute{
			public.NewWebPublic(ctx),
			public.NewLogin(ctx),
			public.NewSignUp(ctx),
			public.NewHome(ctx),
		},
// DB-START
		private: []builders.IAppPrivateRoute{
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
		a.router.Mount(r.Path(), builder.Build())
		lgr.Info(r.Path())
	}

	err := a.appCtx.SM().PrivilegesService().BuildCache()
	if err != nil {
		lgr.Fatal("Error building permission cache", zap.Error(err))
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
