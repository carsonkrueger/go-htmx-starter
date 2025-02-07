package builders

import "github.com/go-chi/chi/v5"

type privateRouteBuilder struct {
	router chi.Router
}

func NewPrivateRouteBuilder(path string) privateRouteBuilder {
	return privateRouteBuilder{
		router: chi.NewRouter(),
	}
}

func (rb *privateRouteBuilder) NewHandle() *privateMethodBuilder {
	return &privateMethodBuilder{
		router: rb.router,
	}
}

func (mb *privateRouteBuilder) Build() chi.Router {
	return mb.router
}
