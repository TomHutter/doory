package actions

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"breadcrumbs": func(b []Breadcrumb) template.HTML {
				if len(b) == 0 {
					return template.HTML("")
				}

				crbs := make([]string, 0)
				for _, br := range b {
					crbs = append(crbs, fmt.Sprintf("<a href=\"%s\">%s</a>", br.Path, br.Name))
				}
				return template.HTML(strings.Join(crbs, " >> "))
			},
			"getParentPath": func(b []Breadcrumb, defaultPath template.HTML) template.HTML {
				if len(b) < 2 {
					return template.HTML("/")
				}

				br := b[len(b)-2]
				return template.HTML(br.Path)
			},
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
		},
	})
}
