package routes

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/go-chi/chi/v5"
)

type HelloWorld struct {
	Ctx *types.AppContext
}

func (r HelloWorld) Path() string {
	return "/helloworld"
}

func (r HelloWorld) Route() chi.Router {
	router := chi.NewRouter()
	router.Get("/", get)
	return router
}

func get(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
