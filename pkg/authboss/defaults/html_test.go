package defaults

import (
	"context"
	"strings"
	"testing"
)

func TestRenderSuccess(t *testing.T) {
	t.Parallel()
	r := NewHTMLRenderer("/auth", "../../../web/templates/authboss/")
	err := r.Load("login.html.tpl")
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

	if !strings.Contains(string(o), "/auth/login") {
		t.Error("expected the url to be rendered out for the form post location")
	}

	o, content, err = r.Render(context.Background(), "mailer/confirm.html.tpl", nil)
	if err != nil {
		t.Error(err)
	}

	if content != "text/html" {
		t.Error("context type not set properly")
	}

	if len(o) == 0 {
		t.Error("it should have rendered a template")
	}
}

func TestRenderFail(t *testing.T) {
	t.Parallel()
	r := NewHTMLRenderer("/auth", "../../../web/templates/authboss/")

	_, _, err := r.Render(context.Background(), "doesntexist....html.tpl", nil)
	if !strings.Contains(err.Error(), "the template doesntexist....html.tpl does not exist") {
		t.Error(err)
	}
}

func TestLoadFail(t *testing.T) {
	t.Parallel()
	r := NewHTMLRenderer("/auth", "../../../web/templates/authboss/")
	err := r.Load("./doesntexist....html.tpl")
	if err == nil {
		t.Error("Expected error due to nonexistent file")
	}
}
