package cmd

import (
	"database/sql"
	"flag"

	"github.com/carsonkrueger/main/internal/cfg"
	"github.com/carsonkrueger/main/internal/logger"
	"github.com/carsonkrueger/main/internal/seeders"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func seed() {
	undo := flag.Bool("undo", false, "-undo=true")
	flag.Parse()
	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg)

	db, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		lgr.Fatal("failed to connect to database", zap.Error(err))
	}

	if *undo {
		lgr.Info("Starting undo...")
		err = seeders.UndoPermissions(db)
	} else {
		lgr.Info("Starting seeds...")
		err = seeders.SeedPermissions(db)
	}
	if err != nil {
		lgr.Fatal("failed to seed database", zap.Error(err))
	}
	lgr.Info("Finished")
}
