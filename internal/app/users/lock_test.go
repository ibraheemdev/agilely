package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
)

type testLockHarness struct {
	users *Users
	e     *engine.Engine

	bodyReader *test.BodyReader
	mailer     *test.Emailer
	redirector *test.Redirector
	renderer   *test.Renderer
	responder  *test.Responder
	session    *test.ClientStateRW
	storer     *test.ServerStorer
}

func testLockSetup() *testLockHarness {
	harness := &testLockHarness{}

	harness.e = engine.New()
	harness.bodyReader = &test.BodyReader{}
	harness.mailer = &test.Emailer{}
	harness.redirector = &test.Redirector{}
	harness.renderer = &test.Renderer{}
	harness.responder = &test.Responder{}
	harness.session = test.NewClientRW()
	harness.storer = test.NewServerStorer()

	harness.e.Config.Authboss.LockAfter = 3
	harness.e.Config.Authboss.LockDuration = time.Hour
	harness.e.Config.Authboss.LockWindow = time.Minute

	harness.e.Core.BodyReader = harness.bodyReader
	harness.e.Core.Logger = test.Logger{}
	harness.e.Core.Mailer = harness.mailer
	harness.e.Core.Redirector = harness.redirector
	harness.e.Core.MailRenderer = harness.renderer
	harness.e.Core.Responder = harness.responder
	harness.e.Core.SessionState = harness.session
	harness.e.Core.Database = harness.storer

	harness.users = NewController(harness.e)

	return harness
}

func TestBeforeAuthAllow(t *testing.T) {
	t.Parallel()

	harness := testLockSetup()

	user := &test.User{
		Email:  "test@test.com",
		Locked: time.Time{},
	}
	harness.storer.Users["test@test.com"] = user

	r := test.Request("GET")
	r = r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))
	w := httptest.NewRecorder()

	handled, err := harness.users.ResetLoginAttempts(w, r, false)
	if err != nil {
		t.Error(err)
	}
	if handled {
		t.Error("it shouldn't have been handled")
	}
}

func TestBeforeAuthDisallow(t *testing.T) {
	t.Parallel()

	harness := testLockSetup()

	user := &test.User{
		Email:  "test@test.com",
		Locked: time.Now().UTC().Add(time.Hour),
	}
	harness.storer.Users["test@test.com"] = user

	r := test.Request("GET")
	r = r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))
	w := httptest.NewRecorder()

	handled, err := harness.users.EnsureNotLocked(w, r, false)
	if err != nil {
		t.Error(err)
	}
	if !handled {
		t.Error("it should have been handled")
	}

	if w.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", w.Code)
	}

	opts := harness.redirector.Options
	if opts.RedirectPath != "/login" {
		t.Error("redir path was wrong:", opts.RedirectPath)
	}

	if len(opts.Failure) == 0 {
		t.Error("expected a failure message")
	}
}

func TestAfterAuthSuccess(t *testing.T) {
	t.Parallel()

	harness := testLockSetup()

	last := time.Now().UTC().Add(-time.Hour)
	user := &test.User{
		Email:        "test@test.com",
		AttemptCount: 45,
		LastAttempt:  last,
	}

	harness.storer.Users["test@test.com"] = user

	r := test.Request("GET")
	r = r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))
	w := httptest.NewRecorder()

	handled, err := harness.users.ResetLoginAttempts(w, r, false)
	if err != nil {
		t.Error(err)
	}
	if handled {
		t.Error("it should never be handled")
	}

	user = harness.storer.Users["test@test.com"]
	if 0 != user.GetAttemptCount() {
		t.Error("attempt count wrong:", user.GetAttemptCount())
	}
	if !last.Before(user.GetLastAttempt()) {
		t.Errorf("last attempt should be more recent, old: %v new: %v", last, user.GetLastAttempt())
	}
}

func TestAfterAuthFailure(t *testing.T) {
	t.Parallel()

	harness := testLockSetup()

	user := &test.User{
		Email: "test@test.com",
	}
	harness.storer.Users["test@test.com"] = user

	if IsLocked(harness.storer.Users["test@test.com"]) {
		t.Error("should not be locked")
	}

	r := test.Request("GET")
	w := httptest.NewRecorder()

	var handled bool
	var err error

	for i := 1; i <= 3; i++ {
		if IsLocked(harness.storer.Users["test@test.com"]) {
			t.Error("should not be locked")
		}

		r := r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))
		handled, err = harness.users.UpdateLockAttempts(w, r, false)
		if err != nil {
			t.Fatal(err)
		}

		if i < 3 {
			if handled {
				t.Errorf("%d) should not be handled until lock occurs", i)
			}

			user := harness.storer.Users["test@test.com"]
			if user.GetAttemptCount() != i {
				t.Errorf("attempt count wrong, want: %d, got: %d", i, user.GetAttemptCount())
			}
			if IsLocked(user) {
				t.Error("should not be locked")
			}
		}
	}

	if !handled {
		t.Error("should have been handled at the end")
	}

	if !IsLocked(harness.storer.Users["test@test.com"]) {
		t.Error("should be locked at the end")
	}

	if w.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", w.Code)
	}

	opts := harness.redirector.Options
	if opts.RedirectPath != "/login" {
		t.Error("redir path was wrong:", opts.RedirectPath)
	}

	if len(opts.Failure) == 0 {
		t.Error("expected a failure message")
	}
}

func TestLock(t *testing.T) {
	t.Parallel()

	harness := testLockSetup()

	user := &test.User{
		Email: "test@test.com",
	}
	harness.storer.Users["test@test.com"] = user

	if IsLocked(harness.storer.Users["test@test.com"]) {
		t.Error("should not be locked")
	}

	if err := harness.users.Lock(context.Background(), "test@test.com"); err != nil {
		t.Error(err)
	}

	if !IsLocked(harness.storer.Users["test@test.com"]) {
		t.Error("should be locked")
	}
}

func TestUnlock(t *testing.T) {
	t.Parallel()

	harness := testLockSetup()

	user := &test.User{
		Email:  "test@test.com",
		Locked: time.Now().UTC().Add(time.Hour),
	}
	harness.storer.Users["test@test.com"] = user

	if !IsLocked(harness.storer.Users["test@test.com"]) {
		t.Error("should be locked")
	}

	if err := harness.users.Unlock(context.Background(), "test@test.com"); err != nil {
		t.Error(err)
	}

	if IsLocked(harness.storer.Users["test@test.com"]) {
		t.Error("should no longer be locked")
	}
}

func TestLockMiddlewareAllow(t *testing.T) {
	t.Parallel()

	e := engine.New()
	u := NewController(e)

	called := false
	server := u.LockMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	user := &test.User{
		Locked: time.Now().UTC().Add(-time.Hour),
	}

	r := test.Request("GET")
	r = r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, r)

	if !called {
		t.Error("The user should have been allowed through")
	}
}

func TestLockMiddlewareDisallow(t *testing.T) {
	t.Parallel()

	e := engine.New()
	redirector := &test.Redirector{}
	e.Core.Logger = test.Logger{}
	e.Core.Redirector = redirector
	u := NewController(e)

	called := false
	server := u.LockMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	user := &test.User{
		Locked: time.Now().UTC().Add(time.Hour),
	}

	r := test.Request("GET")
	r = r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, r)

	if called {
		t.Error("The user should not have been allowed through")
	}
	if redirector.Options.Code != http.StatusTemporaryRedirect {
		t.Error("expected a redirect, but got:", redirector.Options.Code)
	}
	if p := redirector.Options.RedirectPath; p != "/login" {
		t.Error("redirect path wrong:", p)
	}
}
