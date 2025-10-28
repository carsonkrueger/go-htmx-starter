package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/internal/templates/constants"
	"github.com/carsonkrueger/main/internal/templates/ui/layouts"
)

func Tab(req *http.Request, tabModels []layouts.TabModel, selectedTabIndex int, tabContent templ.Component) templ.Component {
	target := req.Header.Get("HX-Target")
	if target == constants.TabContentID {
		return tabContent
	}
	return layouts.Tabs(tabModels, selectedTabIndex)
}
