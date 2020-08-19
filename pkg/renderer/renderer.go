package renderer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path"
	"path/filepath"

	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
)

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

// Renderer :
type Renderer struct {
	// url mount path
	mountPath string

	templates map[string]*template.Template
}

// Render a page
func (t *Renderer) Render(ctx context.Context, page string, data authboss.HTMLData) (output []byte, contentType string, err error) {
	tmpl, ok := t.templates[page]
	if !ok {
		return nil, "", fmt.Errorf("the template %s does not exist", page)
	}
	buf := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to render template for page %s: %w", page, err)
	}
	return buf.Bytes(), "text/html", nil
}

// Load a template directory
func (t *Renderer) Load(layoutsDir, templatesDir string) error {
	layouts, err := filepath.Glob(layoutsDir)
	if err != nil {
		return fmt.Errorf("glob pattern is malformed: %w", err)
	}

	templates, err := filepath.Glob(templatesDir)
	if err != nil {
		return fmt.Errorf("glob pattern is malformed: %w", err)
	}

	funcMap := template.FuncMap{
		"mountpathed": func(location string) string {

			return path.Join(t.mountPath, location)
		},
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}

	mainTemplate, err := template.New("main").Funcs(funcMap).Parse(mainTmpl)
	if err != nil {
		return fmt.Errorf("could not parse main template: %w", err)
	}

	for _, tpl := range templates {
		fileName := filepath.Base(tpl)
		files := append(layouts, tpl)

		t.templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			return fmt.Errorf("template has already been executed: %w", err)
		}

		t.templates[fileName] = template.Must(t.templates[fileName].ParseFiles(files...))
	}
	return nil
}
