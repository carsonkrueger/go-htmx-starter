package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/templates/constants"
	"github.com/carsonkrueger/main/templates/page_layouts"
	"github.com/carsonkrueger/main/util"
)

func Tab(req *http.Request, tabModels []page_layouts.TabModel, selectedTabIndex int, tabContent templ.Component) templ.Component {
	hxRequest := util.IsHxRequest(req)
	target := req.Header.Get("HX-Target")
	if target == constants.TabContentID {
		return tabContent
	}
	content := page_layouts.Tabs(tabModels, selectedTabIndex)
	if target == constants.MainContentID {
		return content
	} else if !hxRequest {
		content = page_layouts.Index(page_layouts.MainPageLayout(content))
	}
	return content
}
