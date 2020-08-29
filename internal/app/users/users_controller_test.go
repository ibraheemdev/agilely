package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
)

func TestUsersInit(t *testing.T) {
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

func testSetupContext() (*Users, *http.Request) {
	e := engine.New()
	e.Core.SessionState = &test.ClientStateRW{
		ClientValues: map[string]string{engine.SessionKey: "george-pid"},
	}

	e.Core.Server = &test.ServerStorer{
		Users: map[string]*test.User{
			"george-pid": {Email: "george-pid", Password: "unreadable"},
		},
	}
	r := httptest.NewRequest("GET", "/", nil)
	w := e.NewResponse(httptest.NewRecorder())

	r, err := e.LoadClientState(w, r)
	if err != nil {
		panic(err)
	}

	return NewController(e), r
}

func testSetupContextCached() (*Users, *test.User, *http.Request) {
	e := engine.New()
	wantUser := &test.User{Email: "george-pid", Password: "unreadable"}
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), CTXKeyPID, "george-pid")
	ctx = context.WithValue(ctx, CTXKeyUser, wantUser)
	req = req.WithContext(ctx)

	return NewController(e), wantUser, req
}

func testSetupContextPanic() *Users {
	e := engine.New()
	e.Core.SessionState = &test.ClientStateRW{
		ClientValues: map[string]string{engine.SessionKey: "george-pid"},
	}
	e.Core.Server = test.NewServerStorer()

	return NewController(e)
}

func TestCurrentUserID(t *testing.T) {
	t.Parallel()

	u, r := testSetupContext()

	id, err := u.CurrentUserID(r)
	if err != nil {
		t.Error(err)
	}

	if id != "george-pid" {
		t.Error("got:", id)
	}
}

func TestCurrentUserIDContext(t *testing.T) {
	t.Parallel()

	u, r := testSetupContext()

	id, err := u.CurrentUserID(r)
	if err != nil {
		t.Error(err)
	}

	if id != "george-pid" {
		t.Error("got:", id)
	}
}

func TestCurrentUserIDP(t *testing.T) {
	t.Parallel()

	e := testSetupContextPanic()
	// Overwrite the setup functions state storer
	e.Core.SessionState = test.NewClientRW()

	defer func() {
		if recover().(error) != engine.ErrUserNotFound {
			t.Failed()
		}
	}()

	_ = e.CurrentUserIDP(httptest.NewRequest("GET", "/", nil))
}

func TestCurrentUser(t *testing.T) {
	t.Parallel()

	u, r := testSetupContext()

	user, err := u.CurrentUser(r)
	if err != nil {
		t.Error(err)
	}

	if got := user.GetPID(); got != "george-pid" {
		t.Error("got:", got)
	}
}

func TestCurrentUserContext(t *testing.T) {
	t.Parallel()

	u, _, r := testSetupContextCached()

	user, err := u.CurrentUser(r)
	if err != nil {
		t.Error(err)
	}

	if got := user.GetPID(); got != "george-pid" {
		t.Errorf("got: %s", got)
	}
}

func TestCurrentUserP(t *testing.T) {
	t.Parallel()

	e := testSetupContextPanic()

	defer func() {
		if recover().(error) != engine.ErrUserNotFound {
			t.Failed()
		}
	}()

	_ = e.CurrentUserP(httptest.NewRequest("GET", "/", nil))
}

func TestLoadCurrentUserID(t *testing.T) {
	t.Parallel()

	u, r := testSetupContext()

	id, err := u.LoadCurrentUserID(&r)
	if err != nil {
		t.Error(err)
	}

	if id != "george-pid" {
		t.Error("got:", id)
	}

	if r.Context().Value(CTXKeyPID).(string) != "george-pid" {
		t.Error("context was not updated in local request")
	}
}

func TestLoadCurrentUserIDContext(t *testing.T) {
	t.Parallel()

	u, _, r := testSetupContextCached()

	pid, err := u.LoadCurrentUserID(&r)
	if err != nil {
		t.Error(err)
	}

	if pid != "george-pid" {
		t.Error("got:", pid)
	}
}

func TestLoadCurrentUserIDP(t *testing.T) {
	t.Parallel()

	e := testSetupContextPanic()

	defer func() {
		if recover().(error) != engine.ErrUserNotFound {
			t.Failed()
		}
	}()

	r := httptest.NewRequest("GET", "/", nil)
	_ = e.LoadCurrentUserIDP(&r)
}

func TestLoadCurrentUser(t *testing.T) {
	t.Parallel()

	u, r := testSetupContext()

	user, err := u.LoadCurrentUser(&r)
	if err != nil {
		t.Error(err)
	}

	if got := user.GetPID(); got != "george-pid" {
		t.Error("got:", got)
	}

	want := user.(*test.User)
	got := r.Context().Value(CTXKeyUser).(*test.User)
	if got != want {
		t.Errorf("users mismatched:\nwant: %#v\ngot: %#v", want, got)
	}
}

func TestLoadCurrentUserContext(t *testing.T) {
	t.Parallel()

	u, wantUser, r := testSetupContextCached()

	user, err := u.LoadCurrentUser(&r)
	if err != nil {
		t.Error(err)
	}

	got := user.(*test.User)
	if got != wantUser {
		t.Errorf("users mismatched:\nwant: %#v\ngot: %#v", wantUser, got)
	}
}

func TestLoadCurrentUserP(t *testing.T) {
	t.Parallel()

	e := testSetupContextPanic()

	defer func() {
		if recover().(error) != engine.ErrUserNotFound {
			t.Failed()
		}
	}()

	r := httptest.NewRequest("GET", "/", nil)
	_ = e.LoadCurrentUserP(&r)
}

func TestCTXKeyString(t *testing.T) {
	t.Parallel()

	if got := CTXKeyUser.String(); got != "users ctx key user" {
		t.Error(got)
	}
}
