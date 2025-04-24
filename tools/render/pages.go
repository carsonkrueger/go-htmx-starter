package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/tools"
)

func PageMainLayout(req *http.Request, page templ.Component) templ.Component {
	hxRequest := tools.IsHxRequest(req)
	target := req.Header.Get("HX-Target")
	if target == pageLayouts.PageLayoutID {
		page = pageLayouts.MainPageLayout(page)
	} else if !hxRequest {
		page = pageLayouts.Index(pageLayouts.MainPageLayout(page))
	}
	return page
}
