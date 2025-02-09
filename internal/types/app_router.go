package types

import (
	"fmt"
	"os"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/go-chi/chi/v5"
)

type RoutePath interface {
	Path() string
}

type PublicRoute interface {
	PublicRoute(r chi.Router)
}

type AppPublicRoute interface {
	SetCtx
	RoutePath
	PublicRoute
}

type PrivateRoute interface {
	PrivateRoute(b *builders.PrivateRouteBuilder)
}

type AppPrivateRoute interface {
	SetCtx
	RoutePath
	PrivateRoute
}

func ReportIfErr(e error, db any) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
		// report error to db
	}
}
