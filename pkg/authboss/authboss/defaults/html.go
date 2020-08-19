package defaults

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path"
	"path/filepath"

	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
)

// HTMLRenderer :
type HTMLRenderer struct {
	// url mount path
	mountPath string

	templatesDir string
	templates    map[string]*template.Template

	// path to layout template. 
	// templates are rendered without
	// a layout if this field is empty
	layout string
}

// NewHTMLRenderer :
func NewHTMLRenderer(mountPath, templatesDir, layout string) *HTMLRenderer {
	return &HTMLRenderer{
		mountPath:    mountPath,
		templates:    make(map[string]*template.Template),
		templatesDir: templatesDir,
		layout:       layout,
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

	if len(r.layout) != 0 {
		name = r.layout
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

	var Layout *template.Template
	if len(r.layout) != 0 {
		l, err := template.New("layout").Funcs(funcMap).ParseFiles(r.layout)
		Layout = l
		if err != nil {
			return fmt.Errorf("could not parse main template: %w", err)
		}
	}

	for _, tpl := range templates {
		filePath := fmt.Sprintf("%s/%s", r.templatesDir, tpl)

		var Files []string = []string{filePath}

		if len(r.layout) != 0 {
			clone, err := Layout.Clone()
			if err != nil {
				return err
			}
			r.templates[tpl] = clone
			Files = append(Files, r.layout)
		}

		template, err := r.templates[tpl].ParseFiles(Files...)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", tpl, err)
		}
		r.templates[tpl] = template
	}
	return nil
}
