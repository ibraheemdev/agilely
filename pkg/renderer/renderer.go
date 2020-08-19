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
	return &HTMLRenderer{
		mountPath:    mountPath,
		templates:    make(map[string]*template.Template),
		templatesDir: templatesDir,
		layoutsDir:   layoutsDir,
	}
}

// NewMailRenderer : Returns a new HTML renderer without a template directory
// because the default mailer templates are standalone
func NewMailRenderer(mountPath, templatesDir string) *HTMLRenderer {
	return &HTMLRenderer{
		mountPath:    mountPath,
		templates:    make(map[string]*template.Template),
		templatesDir: templatesDir,
	}
}

// Render a page
func (r *HTMLRenderer) Render(ctx context.Context, name string, data authboss.HTMLData) (output []byte, contentType string, err error) {
	tmpl, ok := r.templates[name]
	if !ok {
		return nil, "", fmt.Errorf("the template %s does not exist", name)
	}

	buf := &bytes.Buffer{}

	if len(r.layoutsDir) != 0 {
		name = "base"
	}

	err = tmpl.ExecuteTemplate(buf, filepath.Base(name), data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to render template for page %s: %w", name, err)
	}

	return buf.Bytes(), "text/html", nil
}

// Load a template directory
func (r *HTMLRenderer) Load(templates ...string) error {
	funcMap := template.FuncMap{
		"mountpathed": func(location string) string {
			return path.Join(r.mountPath, location)
		},
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}

	for _, tpl := range templates {
		filePath := fmt.Sprintf("%s/%s", r.templatesDir, tpl)

		var Files []string = []string{filePath}

		if len(r.layoutsDir) != 0 {
			layouts, err := filepath.Glob(r.layoutsDir)
			if err != nil {
				return fmt.Errorf("could not parse layouts glob %s, %w", layouts, err)
			}

			mainTemplate, err := template.New("main").Funcs(funcMap).Parse(mainTmpl)
			if err != nil {
				return err
			}

			r.templates[tpl], err = mainTemplate.Clone()

			Files = append(Files, layouts...)
		}

		template, err := r.templates[tpl].ParseFiles(Files...)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", tpl, err)
		}
		r.templates[tpl] = template
	}
	return nil
}
