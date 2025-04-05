package cmd

import (
	"database/sql"
	"flag"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/seeders"
	_ "github.com/lib/pq"
)

func seed() {
	undo := flag.Bool("undo", false, "-undo=true")
	flag.Parse()
	cfg := cfg.LoadConfig()

	db, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}

	lgr := logger.NewLogger(&cfg)

	if *undo {
		lgr.Info("Starting undo...")
		err = seeders.UndoPermissions(db)
	} else {
		lgr.Info("Starting seeds...")
		err = seeders.SeedPermissions(db)
	}
	if err != nil {
		panic(err)
	}
	lgr.Info("Finished")
}
