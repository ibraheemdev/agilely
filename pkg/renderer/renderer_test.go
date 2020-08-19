package renderer

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ibraheemdev/agilely/test"
)

func TestHTMLRenderSuccess(t *testing.T) {
	t.Parallel()
	test.MoveToRoot()

	r := NewHTMLRenderer("/auth", "web/templates/authboss", "web/templates/layout.html.tpl")
	err := r.Load("login.html.tpl", "register.html.tpl")
	if err != nil {
		t.Error(err)
	}

	o, content, err := r.Render(context.Background(), "login.html.tpl", nil)
	fmt.Println(string(o))
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

	if !strings.Contains(string(o), "<!-- Application Layout -->") {
		t.Error("expected the template to be rendered within the layout")
	}
}

func TestMailRenderSuccess(t *testing.T) {
	test.MoveToRoot()
	r := NewMailRenderer("/auth", "web/templates/authboss/mailer")
	err := r.Load("confirm.html.tpl")
	if err != nil {
		t.Error(err)
	}
	o, content, err := r.Render(context.Background(), "confirm.html.tpl", nil)
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
	test.MoveToRoot()
	r := NewHTMLRenderer("/auth", "web/templates/authboss", "web/templates/layout.html.tpl")

	_, _, err := r.Render(context.Background(), "doesntexist....html.tpl", nil)
	if !strings.Contains(err.Error(), "the template doesntexist....html.tpl does not exist") {
		t.Error(err)
	}
}

func TestLoadFail(t *testing.T) {
	t.Parallel()
	test.MoveToRoot()
	r := NewHTMLRenderer("/auth", "web/templates/authboss", "web/templates/layout.html.tpl")
	err := r.Load("doesntexist....html.tpl")
	if err == nil {
		t.Error("Expected error due to nonexistent file")
	}
}
