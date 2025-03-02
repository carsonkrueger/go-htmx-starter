package main

import (
	"github.com/carsonkrueger/main/cfg"
	router "github.com/carsonkrueger/main/internal"
	app_context "github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/logger"
)

func main() {
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)
	ctx := app_context.NewAppContext(lgr, &cfg)
	defer ctx.CleanUp()

	appRouter := router.Setup()
	appRouter.BuildRouter(ctx)
	err := appRouter.Start(cfg)

	if err != nil {
		panic(err)
	}
}
