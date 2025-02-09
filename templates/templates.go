package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	var paths []string
	err := filepath.Walk("templates", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return &Templates{
		templates: template.Must(template.ParseFiles(paths...)),
	}
}
