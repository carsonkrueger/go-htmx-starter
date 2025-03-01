package public_routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-chi/chi/v5"
)

type WebPublic struct {
	types.WithAppContext
}

func (w *WebPublic) Path() string {
	return "/public"
}

func (w *WebPublic) PublicRoute(r chi.Router) {
	r.Get("/{name}", w.get)
}

func (w *WebPublic) get(res http.ResponseWriter, req *http.Request) {
	filename := chi.URLParam(req, "name")
	if filename == "" {
		res.WriteHeader(404)
		return
	}

	dir, _ := os.Getwd()

	filename = fmt.Sprintf("%v/../../assets/%v", dir, filename)

	f, _ := os.Open(filename)

	info, _ := f.Stat()

	contentType := tools.GetMimeType(filename)

	res.Header().Set("Content-Type", contentType)
	res.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	_, _ = f.WriteTo(res)
}
