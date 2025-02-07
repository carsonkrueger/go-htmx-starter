package main

import (
	"fmt"

	router "github.com/carsonkrueger/main/internal"
	app_context "github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/templates"
)

func main() {
	ctx := app_context.AppContext{
		Templates: templates.NewTemplates(),
	}

	appRouter := router.AppRouter{}
	appRouter.Setup(&ctx)
	appRouter.BuildRouter()
	err := appRouter.Start("0.0.0.0", 3000)

	if err != nil {
		fmt.Println(err)
	}
}
