package main

import (
	"database/sql"
	"flag"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/seeders"
	_ "github.com/lib/pq"
)

func main() {
	undo := flag.Bool("undo", false, "Specify if undoing seeds")
	cfg := cfg.LoadConfig()

	db, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}

	if *undo {
		err = seeders.UndoPermissions(db)
	} else {
		err = seeders.SeedPermissions(db)
	}
	if err != nil {
		panic(err)
	}
}
