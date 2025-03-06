package cmd

import (
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/internal"
	"github.com/carsonkrueger/main/internal/services"
	"github.com/carsonkrueger/main/logger"
)

func main() {
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)

	db, err := sql.Open("postgres", cfg.DbUrl())
	// defer db.Close()
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("Database connection is nil")
	}
	sm := services.NewServiceManager(db)

	ctx := internal.NewAppContext(lgr, sm)
	// defer ctx.CleanUp()

	appRouter := Setup()
	appRouter.BuildRouter(ctx)
	err = appRouter.Start(cfg)

	if err != nil {
		panic(err)
	}
}
