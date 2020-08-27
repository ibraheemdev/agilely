package config

import (
	"testing"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
)

func TestRoutes(t *testing.T) {
	e := engine.New()

	router := &test.Router{}
	e.Core.Router = router
	e.Core.ErrorHandler = &test.ErrorHandler{}
	e.Core.ViewRenderer = &test.Renderer{}
	e.Core.MailRenderer = &test.Renderer{}
	e.Core.Server = &test.ServerStorer{}

	if err := Routes(e); err != nil {
		t.Error(err)
	}

	if err := router.HasGets("/polls/:id", "/login", "/confirm", "/recover", "/recover/end", "/register"); err != nil {
		t.Error(err)
	}

	if err := router.HasPuts("/polls/:id"); err != nil {
		t.Error(err)
	}

	if err := router.HasPosts("/polls", "/login", "/recover", "/recover/end", "/register"); err != nil {
		t.Error(err)
	}

	if err := router.HasDeletes("/logout"); err != nil {
		t.Error(err)
	}
}
