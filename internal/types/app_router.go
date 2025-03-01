package types

import (
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
	SetAppContext
	RoutePath
	PublicRoute
}

type PrivateRoute interface {
	PrivateRoute(b *builders.PrivateRouteBuilder)
}

type AppPrivateRoute interface {
	SetAppContext
	RoutePath
	PrivateRoute
}
