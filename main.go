package main

import (
	"fmt"
	"os"

	"github.com/carsonkrueger/main/cmd"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: go run .", "{cmd}")
		os.Exit(1)
	}
	name := os.Args[len(os.Args)-1]
	cmd.Execute(name)
}
