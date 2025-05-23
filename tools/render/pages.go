package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/templates/page_layouts"
	"github.com/carsonkrueger/main/tools"
)

func PageMainLayout(req *http.Request, page templ.Component) templ.Component {
	hxRequest := tools.IsHxRequest(req)
	target := req.Header.Get("HX-Target")
	if target == page_layouts.PageLayoutID {
		page = page_layouts.MainPageLayout(page)
	} else if !hxRequest {
		page = page_layouts.Index(page_layouts.MainPageLayout(page))
	}
	return page
}
