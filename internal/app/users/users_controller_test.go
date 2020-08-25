package users

import (
	"testing"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
)

func TestEngineInit(t *testing.T) {
	t.Parallel()

	e := engine.New()

	router := &test.Router{}
	renderer := &test.Renderer{}
	errHandler := &test.ErrorHandler{}
	server := &test.ServerStorer{}
	mailRenderer := &test.Renderer{}
	e.Config.Core.Router = router
	e.Config.Core.ViewRenderer = renderer
	e.Config.Core.ErrorHandler = errHandler
	e.Config.Core.MailRenderer = mailRenderer
	e.Config.Storage.Server = server

	u := NewController(e)

	if err := u.Init(); err != nil {
		t.Fatal(err)
	}

	if err := renderer.HasLoadedViews(PageLogin, PageRecoverStart, PageRecoverEnd, PageRegister); err != nil {
		t.Error(err)
	}

	if err := router.HasGets("/login", "/confirm", "/recover", "/recover/end", "/register"); err != nil {
		t.Error(err)
	}

	if err := router.HasPosts("/login", "/recover", "/recover/end", "/register"); err != nil {
		t.Error(err)
	}

	if err := router.HasDeletes("/logout"); err != nil {
		t.Error(err)
	}
}
