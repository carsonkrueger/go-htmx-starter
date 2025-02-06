package helloworld

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HWRouter struct {
	test string
}

func (r HWRouter) Path() string {
	return "/helloworld"
}

func (r HWRouter) Route() chi.Router {
	router := chi.NewRouter()
	router.Get("/", get)
	return router
}

func get(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
