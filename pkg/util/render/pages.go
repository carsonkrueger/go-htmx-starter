package render

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/pkg/util"
)

type TargetFunc func(content templ.Component) templ.ComponentFunc
type TemplateTargets map[string]TargetFunc

// templateTargets must have an "index" template
func Layout(ctx context.Context, req *http.Request, w io.Writer, templateTargets TemplateTargets, page templ.Component) error {
	hxRequest := util.IsHxRequest(req)
	target := req.Header.Get("HX-Target")

	if !hxRequest {
		if tt, ok := templateTargets["index"]; ok {
			return tt(page).Render(ctx, w)
		} else {
			panic("no index template")
		}
	}

	if tt, ok := templateTargets[target]; ok {
		return tt(page).Render(ctx, w)
	}

	return page.Render(ctx, w)
}
