package types

import (
	"github.com/go-chi/chi/v5"
)

type AppRouteCtx struct {
	Ctx AppContext
}

type AppRoute interface {
	Path() string
	Route() chi.Router
}
