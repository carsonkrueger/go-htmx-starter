package cmd

import (
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/database/DAO"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/router"
	"github.com/carsonkrueger/main/services"

	_ "github.com/lib/pq"
)

func web() {
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)

// DB-START
	db, err := sql.Open("postgres", cfg.DbUrl())
	defer db.Close()
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("Database connection is nil")
	}

	dm := DAO.NewDAOManager(db)
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
	appRouter.BuildRouter()
	if err := appRouter.Start(cfg); err != nil {
		panic(err)
	}
}
