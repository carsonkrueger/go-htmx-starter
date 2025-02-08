package builders

import "github.com/go-chi/chi/v5"

type PrivateRouteBuilder struct {
	router chi.Router
}

func NewPrivateRouteBuilder() PrivateRouteBuilder {
	return PrivateRouteBuilder{
		router: chi.NewRouter(),
	}
}

func (rb *PrivateRouteBuilder) NewHandle() *privateMethodBuilder {
	return &privateMethodBuilder{
		router: rb.router,
	}
}

func (mb *PrivateRouteBuilder) Build() chi.Router {
	return mb.router
}
