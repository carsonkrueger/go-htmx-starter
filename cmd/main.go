package main

import (
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	router "github.com/carsonkrueger/main/internal"
	app_context "github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/logger"
)

func main() {
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)

	db, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}
	sm := app_context.NewServiceManager(db)

	ctx := app_context.NewAppContext(lgr, sm)
	defer ctx.CleanUp()

	appRouter := router.Setup()
	appRouter.BuildRouter(ctx)
	err = appRouter.Start(cfg)

	if err != nil {
		panic(err)
	}
}
