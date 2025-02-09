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
	types.WithContext
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

	dir, err := os.Getwd()
	types.ReportIfErr(err, nil)

	filename = fmt.Sprintf("%v/web/public/%v", dir, filename)

	f, err := os.Open(filename)
	types.ReportIfErr(err, nil)

	info, err := f.Stat()
	types.ReportIfErr(err, nil)

	contentType := tools.GetMimeType(filename)

	res.Header().Set("Content-Type", contentType)
	res.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	_, err = f.WriteTo(res)
	types.ReportIfErr(err, nil)
}
