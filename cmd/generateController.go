package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

func generateController() {
	controller := flag.String("name", "", "camelCase Name of the controller")
	private := flag.Bool("private", true, "Is a private controller")
	flag.Parse()

	// lower first letter of table name
	controller = tools.Ptr(tools.ToLowerFirst(*controller))

	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg).Named("generateController")

	wd, err := os.Getwd()
	if err != nil {
		lgr.Error("failed to get working directory", zap.Error(err))
		os.Exit(1)
	}

	var filePath string
	if *private {
		filePath = fmt.Sprintf("%s/controllers/private/%s.go", wd, *controller)
	} else {
		filePath = fmt.Sprintf("%s/controllers/public/%s.go", wd, *controller)
	}
	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		lgr.Error("failed to create directory", zap.Error(err))
		os.Exit(1)
	}

	var contents string
	if *private {
		contents = privateControllerFileContents(*controller)
	} else {
		contents = publicControllerFileContents(*controller)
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			lgr.Error("File already exists\n")
			return
		}
		lgr.Error("failed to open file", zap.Error(err))
		return
	}
	io.WriteString(file, contents)
	file.Close()
}

func privateControllerFileContents(name string) string {
	upper := tools.ToUpperFirst(name)
	lower := tools.ToLowerFirst(name)

	return fmt.Sprintf(`package private

import (
	"net/http"

	"github.com/carsonkrueger/main/builders"
	"github.com/carsonkrueger/main/interfaces"
)

const (
	%[2]sGet   = "%[2]sGet"
	%[2]sPost  = "%[2]sPost"
	%[2]sPut = "%[2]sPut"
	%[2]sPatch = "%[2]sPatch"
	%[2]sDelete = "%[2]sDelete"
)

type %[1]s struct {
	interfaces.IAppContext
}

func New%[2]s(ctx interfaces.IAppContext) *%[1]s {
	return &%[1]s{
		IAppContext: ctx,
	}
}

func (um %[1]s) Path() string {
	return "/%[1]s"
}

func (um *%[1]s) PrivateRoute(b *builders.PrivateRouteBuilder) {
	b.NewHandle().Register(builders.GET, "/", um.%[1]sGet).SetPermissionName(%[2]sGet).Build()
	b.NewHandle().Register(builders.POST, "/", um.%[1]sPost).SetPermissionName(%[2]sPost).Build()
	b.NewHandle().Register(builders.PUT, "/", um.%[1]sPut).SetPermissionName(%[2]sPut).Build()
	b.NewHandle().Register(builders.PATCH, "/", um.%[1]sPatch).SetPermissionName(%[2]sPatch).Build()
	b.NewHandle().Register(builders.DELETE, "/", um.%[1]sDelete).SetPermissionName(%[2]sDelete).Build()
}

func (r *%[1]s) %[1]sGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sGet")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sPost(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sPost")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sPut(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sPut")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sPatch(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sPatch")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sDelete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sDelete")
	lgr.Info("Called")
}
`, lower, upper)
}

func publicControllerFileContents(name string) string {
	upper := tools.ToUpperFirst(name)
	lower := tools.ToLowerFirst(name)

	return fmt.Sprintf(`package public

import (
	"net/http"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/go-chi/chi/v5"
)

type %[1]s struct {
	interfaces.IAppContext
}

func New%[2]s(ctx interfaces.IAppContext) *%[1]s {
	return &%[1]s{
		IAppContext: ctx,
	}
}

func (r *%[1]s) Path() string {
	return "/%[1]s"
}

func (r *%[1]s) PublicRoute(router chi.Router) {
	router.Get("/", r.%[1]sGet)
	router.Post("/", r.%[1]sPost)
	router.Put("/", r.%[1]sPut)
	router.Patch("/", r.%[1]sPatch)
	router.Delete("/", r.%[1]sDelete)
}

func (r *%[1]s) %[1]sGet(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sGet")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sPost(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sPost")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sPut(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sPut")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sPatch(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sPatch")
	lgr.Info("Called")
}

func (r *%[1]s) %[1]sDelete(res http.ResponseWriter, req *http.Request) {
	lgr := r.Lgr("%[1]sDelete")
	lgr.Info("Called")
}
`, lower, upper)
}
