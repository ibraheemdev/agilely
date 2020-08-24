package renderer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path"
	"path/filepath"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

// HTMLRenderer :
type HTMLRenderer struct {
	// url mount path
	mountPath string

	templatesDir string
	templates    map[string]*template.Template

	// path to layout template.
	// templates are rendered without
	// a layout if this field is empty
	layoutsDir string
}

// NewHTMLRenderer :
func NewHTMLRenderer(mountPath, templatesDir, layoutsDir string) *HTMLRenderer {
	r := &HTMLRenderer{
		mountPath:    mountPath,
		templates:    make(map[string]*template.Template),
		templatesDir: templatesDir,
		layoutsDir:   layoutsDir,
	}
	r.LoadAll()
	return r
}

// Render a page
func (r *HTMLRenderer) Render(ctx context.Context, name string, data engine.HTMLData) (output []byte, contentType string, err error) {
	tmpl, ok := r.templates[name]
	if !ok {
		return nil, "", fmt.Errorf("the template %s does not exist", name)
	}

	buf := &bytes.Buffer{}

	err = tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to render template for page %s: %w", name, err)
	}

	return buf.Bytes(), "text/html", nil
}

// LoadAll templates in the templates directory
func (r *HTMLRenderer) LoadAll() error {
	funcMap := template.FuncMap{
		"mountpathed": func(location string) string {
			return path.Join(r.mountPath, location)
		},
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}

	templates, err := filepath.Glob(r.templatesDir)
	if err != nil {
		return fmt.Errorf("could not parse templates glob %s, %w", templates, err)
	}

	for _, tpl := range templates {
		layouts, err := filepath.Glob(r.layoutsDir)
		if err != nil {
			return fmt.Errorf("could not parse layouts glob %s, %w", layouts, err)
		}

		mainTemplate, err := template.New("main").Funcs(funcMap).Parse(mainTmpl)
		if err != nil {
			return err
		}

		r.templates[filepath.Base(tpl)], err = mainTemplate.Clone()

		files := append(layouts, tpl)
		r.templates[filepath.Base(tpl)], err = r.templates[filepath.Base(tpl)].ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", tpl, err)
		}
	}
	return nil
}

// Load : Checks to see if a particular templates has been loaded
func (r *HTMLRenderer) Load(templates ...string) error {
	for _, tpl := range templates {
		_, ok := r.templates[tpl]
		if !ok {
			return fmt.Errorf("failed to load template: %s", tpl)
		}
	}
	return nil
}
