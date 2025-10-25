package router

import (
	gctx "context"
	"errors"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/cfg"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/controllers/private"
	"github.com/carsonkrueger/main/internal/controllers/public"
	"github.com/carsonkrueger/main/internal/middlewares"

	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

type AppRouter struct {
	public  []builders.AppPublicRoute
	private []builders.AppPrivateRoute
	addr    string
	router  chi.Router
	appCtx  context.AppContext
}

func NewAppRouter(appCtx context.AppContext) AppRouter {
	return AppRouter{
		appCtx: appCtx,
		public: []builders.AppPublicRoute{
			public.NewLogin(appCtx),
			public.NewSignUp(appCtx),
			public.NewWebPublic(appCtx),
			public.NewHome(appCtx),
			public.NewCart(appCtx),
			// INSERT PUBLIC
		},
		private: []builders.AppPrivateRoute{
			private.NewUserManagement(appCtx),
			private.NewPrivileges(appCtx),
			private.NewRoles(appCtx),
			private.NewRolesPrivileges(appCtx),
			private.NewProducts(appCtx),
			// INSERT PRIVATE
		},
	}
}

func (a *AppRouter) BuildRouter(ctx gctx.Context) {
	lgr := a.appCtx.Lgr("BuildRouter")
	a.router = chi.NewRouter()
	a.router = a.router.With(middlewares.Recover(a.appCtx))

	a.router = a.router.Group(func(g chi.Router) {
		// do not enforce auth
		g.Use(middlewares.Auth(a.appCtx, false))
		for _, r := range a.public {
			router := chi.NewRouter()
			r.PublicRoute(router)
			g.Mount(r.Path(), router)
			lgr.Info(r.Path())
		}
	})

	a.router = a.router.Group(func(g chi.Router) {
		// enforce auth middleware
		g.Use(middlewares.Auth(a.appCtx, true))
		for _, r := range a.private {
			builder := builders.NewPrivateRouteBuilder(a.appCtx)
			r.PrivateRoute(ctx, &builder)
			g.Mount(r.Path(), builder.RawRouter())
			lgr.Info(r.Path())
		}
	})
}

func (a *AppRouter) Start(cfg cfg.Config) error {
	if a.router == nil {
		return errors.New("AppRouter has no router. Did you forget to call BuildRouter?")
	}

	a.addr = fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	fmt.Printf("\nListening on http://%s\n", a.addr)
	return http.ListenAndServe(a.addr, a.router)
}
