package engine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func loadClientStateP(e *Engine, w http.ResponseWriter, r *http.Request) *http.Request {
	r, err := e.LoadClientState(w, r)
	if err != nil {
		panic(err)
	}

	return r
}

func testSetupContext() (*Engine, *http.Request) {
	e := New()
	e.Core.SessionState = newMockClientStateRW(SessionKey, "george-pid")
	e.Core.Server = &mockServerStorer{
		Users: map[string]*mockUser{
			"george-pid": {Email: "george-pid", Password: "unreadable"},
		},
	}
	r := httptest.NewRequest("GET", "/", nil)
	w := e.NewResponse(httptest.NewRecorder())
	r = loadClientStateP(e, w, r)

	return e, r
}

func testSetupContextCached() (*Engine, *mockUser, *http.Request) {
	e := New()
	wantUser := &mockUser{Email: "george-pid", Password: "unreadable"}
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), CTXKeyPID, "george-pid")
	ctx = context.WithValue(ctx, CTXKeyUser, wantUser)
	req = req.WithContext(ctx)

	return e, wantUser, req
}

func testSetupContextPanic() *Engine {
	e := New()
	e.Core.SessionState = newMockClientStateRW(SessionKey, "george-pid")
	e.Core.Server = &mockServerStorer{}

	return e
}

func TestCurrentUserID(t *testing.T) {
	t.Parallel()

	e, r := testSetupContext()

	id, err := e.CurrentUserID(r)
	if err != nil {
		t.Error(err)
	}

	if id != "george-pid" {
		t.Error("got:", id)
	}
}

func TestCurrentUserIDContext(t *testing.T) {
	t.Parallel()

	e, r := testSetupContext()

	id, err := e.CurrentUserID(r)
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
	e.Core.SessionState = newMockClientStateRW()

	defer func() {
		if recover().(error) != ErrUserNotFound {
			t.Failed()
		}
	}()

	_ = e.CurrentUserIDP(httptest.NewRequest("GET", "/", nil))
}

func TestCurrentUser(t *testing.T) {
	t.Parallel()

	e, r := testSetupContext()

	user, err := e.CurrentUser(r)
	if err != nil {
		t.Error(err)
	}

	if got := user.GetPID(); got != "george-pid" {
		t.Error("got:", got)
	}
}

func TestCurrentUserContext(t *testing.T) {
	t.Parallel()

	e, _, r := testSetupContextCached()

	user, err := e.CurrentUser(r)
	if err != nil {
		t.Error(err)
	}

	if got := user.GetPID(); got != "george-pid" {
		t.Error("got:", got)
	}
}

func TestCurrentUserP(t *testing.T) {
	t.Parallel()

	e := testSetupContextPanic()

	defer func() {
		if recover().(error) != ErrUserNotFound {
			t.Failed()
		}
	}()

	_ = e.CurrentUserP(httptest.NewRequest("GET", "/", nil))
}

func TestLoadCurrentUserID(t *testing.T) {
	t.Parallel()

	e, r := testSetupContext()

	id, err := e.LoadCurrentUserID(&r)
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

	e, _, r := testSetupContextCached()

	pid, err := e.LoadCurrentUserID(&r)
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
		if recover().(error) != ErrUserNotFound {
			t.Failed()
		}
	}()

	r := httptest.NewRequest("GET", "/", nil)
	_ = e.LoadCurrentUserIDP(&r)
}

func TestLoadCurrentUser(t *testing.T) {
	t.Parallel()

	e, r := testSetupContext()

	user, err := e.LoadCurrentUser(&r)
	if err != nil {
		t.Error(err)
	}

	if got := user.GetPID(); got != "george-pid" {
		t.Error("got:", got)
	}

	want := user.(*mockUser)
	got := r.Context().Value(CTXKeyUser).(*mockUser)
	if got != want {
		t.Errorf("users mismatched:\nwant: %#v\ngot: %#v", want, got)
	}
}

func TestLoadCurrentUserContext(t *testing.T) {
	t.Parallel()

	e, wantUser, r := testSetupContextCached()

	user, err := e.LoadCurrentUser(&r)
	if err != nil {
		t.Error(err)
	}

	got := user.(*mockUser)
	if got != wantUser {
		t.Errorf("users mismatched:\nwant: %#v\ngot: %#v", wantUser, got)
	}
}

func TestLoadCurrentUserP(t *testing.T) {
	t.Parallel()

	e := testSetupContextPanic()

	defer func() {
		if recover().(error) != ErrUserNotFound {
			t.Failed()
		}
	}()

	r := httptest.NewRequest("GET", "/", nil)
	_ = e.LoadCurrentUserP(&r)
}

func TestCTXKeyString(t *testing.T) {
	t.Parallel()

	if got := CTXKeyPID.String(); got != "engine ctx key pid" {
		t.Error(got)
	}
}
