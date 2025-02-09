package public_routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/templates/layouts"
	"github.com/go-chi/chi/v5"
)

type HelloWorld struct {
	types.WithContext
}

func (r *HelloWorld) Path() string {
	return "/"
}

func (hw *HelloWorld) PublicRoute(r chi.Router) {
	r.Get("/", hw.index)
	r.Get("/home", hw.home)
	r.Get("/about", hw.about)
}

func (hw *HelloWorld) index(res http.ResponseWriter, req *http.Request) {
	home := layouts.Index(layouts.Home())
	err := home.Render(context.Background(), res)
	fmt.Println(err)
}

func (hw *HelloWorld) home(res http.ResponseWriter, req *http.Request) {
	err := hw.GetCtx().Templates.Render(res, "home.html", nil)
	fmt.Println(err)
}

func (hw *HelloWorld) about(res http.ResponseWriter, req *http.Request) {
	err := hw.GetCtx().Templates.Render(res, "about.html", nil)
	fmt.Println(err)
}
