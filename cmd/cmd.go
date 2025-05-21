package cmd

import (
	"fmt"
)

func Execute(cmd string) {
	switch cmd {
	case "web":
		web()
		// DB-START
	case "seed":
		seed()
	case "genDAO":
		generateDAO()
		// DB-END
	case "genController":
		generateController()
	case "genService":
		generateService()
	default:
		panic(fmt.Sprintf("Invalid cmd: %s", cmd))
	}
}
