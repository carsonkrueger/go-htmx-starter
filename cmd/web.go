package cmd

import (
	gctx "context"
	"database/sql"

	"github.com/carsonkrueger/main/internal/cfg"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/database/dao"
	"github.com/carsonkrueger/main/internal/logger"
	"github.com/carsonkrueger/main/internal/router"
	"github.com/carsonkrueger/main/internal/services"

	_ "github.com/lib/pq"
)

func web() {
	ctx := gctx.Background()
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)
	defer lgr.Sync()

	db, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if db == nil {
		panic("Database connection is nil")
	}
	ctx = context.WithDB(ctx, db)
	ctx = context.WithLogger(ctx, lgr)

	dm := dao.NewDAOManager()
	sm := services.NewServiceManager(nil)
	appCtx := context.NewAppContext(
		sm,
		dm,
	)
	sm.AppContext = appCtx

	appRouter := router.NewAppRouter(appCtx)
	appRouter.BuildRouter(ctx)
	if err := appRouter.Start(cfg); err != nil {
		panic(err)
	}
}
