package main

import (
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/router"
	"github.com/carsonkrueger/main/services"

	_ "github.com/lib/pq"
)

func main() {
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)

	db, err := sql.Open("postgres", cfg.DbUrl())
	defer db.Close()
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("Database connection is nil")
	}
	sm := services.NewServiceManager(db)

	ctx := context.NewAppContext(lgr, sm)
	defer ctx.CleanUp()

	appRouter := router.Setup()
	appRouter.BuildRouter(ctx)
	err = appRouter.Start(cfg)

	if err != nil {
		panic(err)
	}
}
