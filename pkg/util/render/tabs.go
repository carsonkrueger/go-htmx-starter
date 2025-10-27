package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/internal/templates/constants"
	"github.com/carsonkrueger/main/internal/templates/ui/layouts"
	"github.com/carsonkrueger/main/pkg/util"
)

func Tab(req *http.Request, tabModels []layouts.TabModel, selectedTabIndex int, tabContent templ.Component) templ.Component {
	hxRequest := util.IsHxRequest(req)
	target := req.Header.Get("HX-Target")
	if target == constants.TabContentID {
		return tabContent
	}
	content := layouts.Tabs(tabModels, selectedTabIndex)
	if target == constants.MainContentID {
		return content
	} else if !hxRequest {
		// content = layouts.Index(layouts.MainPageLayout(content))
		content = layouts.Index()
	}
	return content
}
