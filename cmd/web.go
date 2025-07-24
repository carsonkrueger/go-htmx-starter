package cmd

import (
	gctx "context"
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/database/dao"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/router"
	"github.com/carsonkrueger/main/services"

	_ "github.com/lib/pq"
)

func web() {
	ctx := gctx.Background()
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)

	// DB-START
	db, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if db == nil {
		panic("Database connection is nil")
	}
	ctx = context.WithDB(ctx, db)

	dm := dao.NewDAOManager()
	// DB-END
	sm := services.NewServiceManager(nil)
	appCtx := context.NewAppContext(
		lgr,
		sm,
		// DB-START
		dm,
		db,
		// DB-END
	)
	sm.SetAppContext(appCtx)
	defer appCtx.CleanUp()

	appRouter := router.NewAppRouter(appCtx)
	appRouter.BuildRouter(ctx)
	if err := appRouter.Start(cfg); err != nil {
		panic(err)
	}
}
