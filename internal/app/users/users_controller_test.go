package users

import (
	"reflect"
	"testing"
	"unsafe"

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
	e.Core.Router = router
	e.Core.ViewRenderer = renderer
	e.Core.ErrorHandler = errHandler
	e.Core.MailRenderer = mailRenderer
	e.Core.Server = server

	u := NewController(e)

	if err := u.Init(); err != nil {
		t.Fatal(err)
	}

	if err := renderer.HasLoadedViews(PageLogin, PageRecoverStart, PageRecoverEnd, PageRegister); err != nil {
		t.Error(err)
	}

	events := reflect.ValueOf(e.AuthEvents).Elem()

	a := events.FieldByName("after")
	after := reflect.NewAt(a.Type(), unsafe.Pointer(a.UnsafeAddr())).Elem().Interface().(map[engine.AuthEvent][]engine.AuthEventHandler)

	if len(after) != 5 {
		t.Errorf("expected 1 event, got %d", len(after))
	}

	if len(after[engine.EventRegister]) != 1 {
		t.Errorf("expected 1 event, got %d", len(after[engine.EventRegister]))
	}

	if len(after[engine.EventAuth]) != 3 {
		t.Errorf("expected 3 events, got %d", len(after[engine.EventAuth]))
	}

	if len(after[engine.EventOAuth2]) != 1 {
		t.Errorf("expected 1 event, got %d", len(after[engine.EventOAuth2]))
	}

	if len(after[engine.EventAuthFail]) != 1 {
		t.Errorf("expected 1 event, got %d", len(after[engine.EventAuthFail]))
	}

	if len(after[engine.EventPasswordReset]) != 1 {
		t.Errorf("expected 1 event, got %d", len(after[engine.EventPasswordReset]))
	}

	b := events.FieldByName("before")
	before := reflect.NewAt(b.Type(), unsafe.Pointer(b.UnsafeAddr())).Elem().Interface().(map[engine.AuthEvent][]engine.AuthEventHandler)

	if len(before) != 2 {
		t.Errorf("expected 1 event, got %d", len(before))
	}

	if len(before[engine.EventAuth]) != 2 {
		t.Errorf("expected 2 events, got %d", len(before[engine.EventAuth]))
	}

	if len(before[engine.EventOAuth2]) != 1 {
		t.Errorf("expected 1 event, got %d", len(before[engine.EventOAuth2]))
	}
}
