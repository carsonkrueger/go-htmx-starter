package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/templates/constants"
	"github.com/carsonkrueger/main/templates/page_layouts"
	"github.com/carsonkrueger/main/util"
)

func PageMainLayout(req *http.Request, page templ.Component) templ.Component {
	hxRequest := util.IsHxRequest(req)
	target := req.Header.Get("HX-Target")
	if target == constants.PageLayoutID {
		page = page_layouts.MainPageLayout(page)
	} else if !hxRequest {
		page = page_layouts.Index(page_layouts.MainPageLayout(page))
	}
	return page
}
