package template

import (
	"io"

	"text/template"
)

const (
	PrivateController = "private-controller"
	PublicController  = "public-controller"
)

var templates = map[string]string{
	PrivateController: privateController,
}

func ExecuteTemplate(wr io.Writer, name string, model any) {
	t, err := template.New(name).Parse(templates[name])
	if err != nil {
		panic(err)
	}
	if err = t.ExecuteTemplate(wr, name, model); err != nil {
		panic(err)
	}
}

var privateController = `package private

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/context"
)

type {{ .NameLower }} struct {
	*context.AppContext
}

func New{{ .Name }}(ctx *context.AppContext) *{{ .NameLower }} {
	return &{{ .NameLower }}{
		AppContext: ctx,
	}
}

func (r {{ .NameLower }}) Path() string {
	return "/{{ .NameLower }}"
}

func (r *{{ .NameLower }}) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/", r.{{ .NameLower }}Get).Build(ctx)
	b.NewHandler().Register(http.MethodPost, "/", r.{{ .NameLower }}Post).Build(ctx)
	b.NewHandler().Register(http.MethodPut, "/", r.{{ .NameLower }}Put).Build(ctx)
	b.NewHandler().Register(http.MethodPatch, "/", r.{{ .NameLower }}Patch).Build(ctx)
	b.NewHandler().Register(http.MethodDelete, "/", r.{{ .NameLower }}Delete).Build(ctx)
}

func (r *{{ .NameLower }}) {{ .NameLower }}Get(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Get")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Post(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Post")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Put(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Put")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Patch(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Patch")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Delete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Delete")
	lgr.Info("Called")
}
`

var publicController = `package public

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/go-chi/chi/v5"
)

type {{ .NameLower }} struct {
	context.AppContext
}

func New{{ .Name }}(ctx context.AppContext) *{{ .NameLower }} {
	return &{{ .NameLower }}{
		AppContext: ctx,
	}
}

func (r *{{ .NameLower }}) Path() string {
	return "/{{ .NameLower }}"
}

func (r *{{ .NameLower }}) PublicRoute(router chi.Router) {
	router.Get("/", r.{{ .NameLower }}Get)
	router.Post("/", r.{{ .NameLower }}Post)
	router.Put("/", r.{{ .NameLower }}Put)
	router.Patch("/", r.{{ .NameLower }}Patch)
	router.Delete("/", r.{{ .NameLower }}Delete)
}

func (r *{{ .NameLower }}) {{ .NameLower }}Get(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Get")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Post(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Post")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Put(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Put")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Patch(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Patch")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Delete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("{{ .NameLower }}Delete")
	lgr.Info("Called")
}
`
