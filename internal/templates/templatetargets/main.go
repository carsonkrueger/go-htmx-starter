package templatetargets

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/internal/templates/constants"
	"github.com/carsonkrueger/main/internal/templates/ui/layouts"
	"github.com/carsonkrueger/main/pkg/util/render"
)

var Main = render.TemplateTargets{
	"index": func(content templ.Component) templ.ComponentFunc {
		return func(ctx context.Context, w io.Writer) error {
			return layouts.Index().Render(templ.WithChildren(ctx, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				return layouts.MainPageLayout().Render(templ.WithChildren(ctx, content), w)
			})), w)
		}
	},
	constants.PageLayoutID: func(content templ.Component) templ.ComponentFunc {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			return layouts.MainPageLayout().Render(templ.WithChildren(ctx, content), w)
		})
	},
}
