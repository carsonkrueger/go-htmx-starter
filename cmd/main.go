package main

import (
	"fmt"

	router "github.com/carsonkrueger/main/internal"
	app_context "github.com/carsonkrueger/main/internal/types"
)

func main() {
	ctx := app_context.AppContext{}

	appRouter := router.AppRouter{}
	appRouter.Setup()
	appRouter.BuildRouter(&ctx)
	err := appRouter.Start("0.0.0.0", 3000)

	if err != nil {
		fmt.Println(err)
	}
}
