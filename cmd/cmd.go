package cmd

import (
	"fmt"
	"os"
)

var cmds = map[string]func(){
	"web":           web,
	"seed":          seed,
	"genDAO":        generateDAO,
	"genController": generateController,
	"genService":    generateService,
}

func Execute(cmd string) {
	if f, ok := cmds[cmd]; ok {
		f()
		os.Exit(0)
	}
	panic(fmt.Sprintf("Invalid cmd: %s", cmd))
}
