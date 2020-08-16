package renderer

import (
	"context"
	"html/template"
	"strings"
	"testing"
)

func TestRenderSuccess(t *testing.T) {
	r := &Renderer{
		templates: make(map[string]*template.Template),
		mountPath: "/",
	}
	err := r.Load("../../web/templates/layouts/*.tpl", "../../web/templates/authboss/*.tpl")
	if err != nil {
		t.Error(err)
	}

	o, content, err := r.Render(context.Background(), "login.html.tpl", nil)
	if err != nil {
		t.Error(err)
	}

	if content != "text/html" {
		t.Error("context type not set properly")
	}

	if len(o) == 0 {
		t.Error("it should have rendered a template")
	}

	if !strings.Contains(string(o), "/login") {
		t.Error("expected the url to be rendered out for the form post location")
	}

	if !strings.Contains(string(o), "<!-- Application Layout -->") {
		t.Error("expected the page to be rendered within a layout")
	}
}

func TestRenderFail(t *testing.T) {
	r := &Renderer{
		templates: make(map[string]*template.Template),
		mountPath: "/",
	}
	err := r.Load("../../web/templates/layouts/*.tpl", "../../web/templates/authboss/*.tpl")
	if err != nil {
		t.Error(err)
	}

	_, _, err = r.Render(context.Background(), "doesntexist....html.tpl", nil)
	if !strings.Contains(err.Error(), "the template doesntexist....html.tpl does not exist") {
		t.Error(err)
	}
}
func TestLoadFail(t *testing.T) {
	r := &Renderer{
		templates: make(map[string]*template.Template),
	}
	mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}asdasd`
	err := r.Load("../../web/templates/layouts/*.tpl", "../../web/templates/authboss/*.tpl")
	if err == nil {
		t.Error("expected error due to malformed main template")
	}
}
