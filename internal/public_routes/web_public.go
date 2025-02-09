package public_routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/carsonkrueger/main/internal/types"
	"github.com/carsonkrueger/main/pkg"
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
	fmt.Printf("%s\n", filename)
	if filename == "" {
		return
	}
	filename = fmt.Sprintf("/home/carson/Repos/go-test/web/public/%v", filename)

	f, err := os.Open(filename)
	types.ReportIfErr(err, nil)

	info, e := f.Stat()
	types.ReportIfErr(e, nil)

	contentType := pkg.GetMimeType(filename)
	fmt.Println(contentType)

	res.Header().Set("Content-Type", contentType)
	res.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	_, err3 := f.WriteTo(res)
	types.ReportIfErr(err3, nil)
}
