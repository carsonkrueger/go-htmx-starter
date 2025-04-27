package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/templates/pageLayouts"
	"github.com/carsonkrueger/main/tools"
)

func Tab(req *http.Request, tabModels []pageLayouts.TabModel, selectedTabIndex int, tabContent templ.Component) templ.Component {
	hxRequest := tools.IsHxRequest(req)
	target := req.Header.Get("HX-Target")
	if target == pageLayouts.TabContentID {
		return tabContent
	}
	content := pageLayouts.Tabs(tabModels, selectedTabIndex)
	if target == pageLayouts.MainContentID {
		return content
	} else if !hxRequest {
		content = pageLayouts.Index(pageLayouts.MainPageLayout(content))
	}
	return content
}
