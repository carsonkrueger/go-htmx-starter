package cmd

import (
	"fmt"
)

func Execute(cmd string) {
	switch cmd {
	case "seed":
		seed()
	case "web":
		web()
	case "genDAO":
		generateDAO()
	case "genController":
		generateController()
	default:
		panic(fmt.Sprintf("Invalid cmd: %s", cmd))
	}
}
