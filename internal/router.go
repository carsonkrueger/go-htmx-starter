package router

import (
	"fmt"
	"net/http"

	helloworld "github.com/carsonkrueger/main/internal/routes/hello_world"
	"github.com/go-chi/chi/v5"
)

type MyRoute interface {
	Path() string
	Route() chi.Router
}

func SetupRouter() {
	publicRoutes := []MyRoute{new(helloworld.HWRouter)}
	privateRoutes := []MyRoute{}

	router(publicRoutes, privateRoutes)
}

func router(publicRoutes []MyRoute, privateRoutes []MyRoute) {
	router := chi.NewRouter()

	fmt.Println("Creating public routes")
	for _, r := range publicRoutes {
		fmt.Printf("Registered %v\n", r.Path())
		router.Mount(r.Path(), r.Route())
	}

	// enforce auth mw
	// router.Use()

	fmt.Println("Creating private routes")
	for _, r := range privateRoutes {
		fmt.Printf("Registered %v\n", r.Path())
		router.Mount(r.Path(), r.Route())
	}

	addr := fmt.Sprintf("%v:%v", "0.0.0.0", 3000)
	fmt.Printf("Listening on %v\n", addr)
	http.ListenAndServe(addr, router)
}
