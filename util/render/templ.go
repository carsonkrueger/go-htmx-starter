package render

import "github.com/a-h/templ"

func SafeTemplJoin(components ...templ.Component) templ.Component {
	var templs []templ.Component
	for _, component := range components {
		if component != nil {
			templs = append(templs, component)
		}
	}
	return templ.Join(templs...)
}
