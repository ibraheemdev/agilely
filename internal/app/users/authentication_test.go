package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
)

func TestAuthGet(t *testing.T) {
	t.Parallel()

	e := engine.New()
	responder := &test.Responder{}
	e.Core.Responder = responder

	u := NewController(e)

	r := test.Request("GET")
	r.URL.RawQuery = "redir=/redirectpage"
	if err := u.LoginGet(nil, r); err != nil {
		t.Error(err)
	}

	if responder.Page != PageLogin {
		t.Error("wanted login page, got:", responder.Page)
	}

	if responder.Status != http.StatusOK {
		t.Error("wanted ok status, got:", responder.Status)
	}

	if got := responder.Data[engine.FormValueRedirect]; got != "/redirectpage" {
		t.Error("redirect page was wrong:", got)
	}
}

type testHarness struct {
	users *Users
	e     *engine.Engine

	bodyReader *test.BodyReader
	responder  *test.Responder
	redirector *test.Redirector
	session    *test.ClientStateRW
	storer     *test.ServerStorer
}

func testSetup() *testHarness {
	harness := &testHarness{}

	harness.e = engine.New()
	harness.bodyReader = &test.BodyReader{}
	harness.redirector = &test.Redirector{}
	harness.responder = &test.Responder{}
	harness.session = test.NewClientRW()
	harness.storer = test.NewServerStorer()

	harness.e.Core.BodyReader = harness.bodyReader
	harness.e.Core.Logger = test.Logger{}
	harness.e.Core.Responder = harness.responder
	harness.e.Core.Redirector = harness.redirector
	harness.e.Core.SessionState = harness.session
	harness.e.Core.Server = harness.storer

	harness.users = NewController(harness.e)

	return harness
}

func TestAuthPostSuccess(t *testing.T) {
	t.Parallel()

	setupMore := func(h *testHarness) *testHarness {
		h.bodyReader.Return = test.Values{
			PID:      "test@test.com",
			Password: "hello world",
		}
		h.storer.Users["test@test.com"] = &test.User{
			Email:    "test@test.com",
			Password: "$2a$10$IlfnqVyDZ6c1L.kaA/q3bu1nkAC6KukNUsizvlzay1pZPXnX2C9Ji", // hello world
		}
		h.session.ClientValues[engine.SessionHalfAuthKey] = "true"

		return h
	}

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		h := setupMore(testSetup())

		var beforeCalled, afterCalled bool
		var beforeHasValues, afterHasValues bool
		h.e.AuthEvents.Before(engine.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			beforeCalled = true
			beforeHasValues = r.Context().Value(engine.CTXKeyValues) != nil
			return false, nil
		})
		h.e.AuthEvents.After(engine.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			afterCalled = true
			afterHasValues = r.Context().Value(engine.CTXKeyValues) != nil
			return false, nil
		})

		r := test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.e.NewResponse(resp)

		if err := h.users.LoginPost(w, r); err != nil {
			t.Error(err)
		}

		if resp.Code != http.StatusTemporaryRedirect {
			t.Error("code was wrong:", resp.Code)
		}
		if h.redirector.Options.RedirectPath != "/" {
			t.Error("redirect path was wrong:", h.redirector.Options.RedirectPath)
		}

		if _, ok := h.session.ClientValues[engine.SessionHalfAuthKey]; ok {
			t.Error("half auth should have been deleted")
		}
		if pid := h.session.ClientValues[engine.SessionKey]; pid != "test@test.com" {
			t.Error("pid was wrong:", pid)
		}

		if !beforeCalled {
			t.Error("before should have been called")
		}
		if !afterCalled {
			t.Error("after should have been called")
		}
		if !beforeHasValues {
			t.Error("before callback should have access to values")
		}
		if !afterHasValues {
			t.Error("after callback should have access to values")
		}
	})

	t.Run("handledBefore", func(t *testing.T) {
		t.Parallel()
		h := setupMore(testSetup())

		var beforeCalled bool
		h.e.AuthEvents.Before(engine.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			w.WriteHeader(http.StatusTeapot)
			beforeCalled = true
			return true, nil
		})

		r := test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.e.NewResponse(resp)

		if err := h.users.LoginPost(w, r); err != nil {
			t.Error(err)
		}

		if h.responder.Status != 0 {
			t.Error("a status should never have been sent back")
		}
		if _, ok := h.session.ClientValues[engine.SessionKey]; ok {
			t.Error("session key should not have been set")
		}

		if !beforeCalled {
			t.Error("before should have been called")
		}
		if resp.Code != http.StatusTeapot {
			t.Error("should have left the response alone once teapot was sent")
		}
	})

	t.Run("handledAfter", func(t *testing.T) {
		t.Parallel()
		h := setupMore(testSetup())

		var afterCalled bool
		h.e.AuthEvents.After(engine.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			w.WriteHeader(http.StatusTeapot)
			afterCalled = true
			return true, nil
		})

		r := test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.e.NewResponse(resp)

		if err := h.users.LoginPost(w, r); err != nil {
			t.Error(err)
		}

		if h.responder.Status != 0 {
			t.Error("a status should never have been sent back")
		}
		if _, ok := h.session.ClientValues[engine.SessionKey]; !ok {
			t.Error("session key should have been set")
		}

		if !afterCalled {
			t.Error("after should have been called")
		}
		if resp.Code != http.StatusTeapot {
			t.Error("should have left the response alone once teapot was sent")
		}
	})
}

func TestAuthPostBadPassword(t *testing.T) {
	t.Parallel()

	setupMore := func(h *testHarness) *testHarness {
		h.bodyReader.Return = test.Values{
			PID:      "test@test.com",
			Password: "world hello",
		}
		h.storer.Users["test@test.com"] = &test.User{
			Email:    "test@test.com",
			Password: "$2a$10$IlfnqVyDZ6c1L.kaA/q3bu1nkAC6KukNUsizvlzay1pZPXnX2C9Ji", // hello world
		}

		return h
	}

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		h := setupMore(testSetup())

		r := test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.e.NewResponse(resp)

		var afterCalled bool
		h.e.AuthEvents.After(engine.EventAuthFail, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			afterCalled = true
			return false, nil
		})

		if err := h.users.LoginPost(w, r); err != nil {
			t.Error(err)
		}

		if resp.Code != 200 {
			t.Error("wanted a 200:", resp.Code)
		}

		if h.responder.Data[engine.DataErr] != "Invalid Credentials" {
			t.Error("wrong error:", h.responder.Data)
		}

		if _, ok := h.session.ClientValues[engine.SessionKey]; ok {
			t.Error("user should not be logged in")
		}

		if !afterCalled {
			t.Error("after should have been called")
		}
	})

	t.Run("handledAfter", func(t *testing.T) {
		t.Parallel()
		h := setupMore(testSetup())

		r := test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.e.NewResponse(resp)

		var afterCalled bool
		h.e.AuthEvents.After(engine.EventAuthFail, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			w.WriteHeader(http.StatusTeapot)
			afterCalled = true
			return true, nil
		})

		if err := h.users.LoginPost(w, r); err != nil {
			t.Error(err)
		}

		if h.responder.Status != 0 {
			t.Error("responder should not have been called to give a status")
		}
		if _, ok := h.session.ClientValues[engine.SessionKey]; ok {
			t.Error("user should not be logged in")
		}

		if !afterCalled {
			t.Error("after should have been called")
		}
		if resp.Code != http.StatusTeapot {
			t.Error("should have left the response alone once teapot was sent")
		}
	})
}

func TestAuthPostUserNotFound(t *testing.T) {
	t.Parallel()

	harness := testSetup()
	harness.bodyReader.Return = test.Values{
		PID:      "test@test.com",
		Password: "world hello",
	}

	r := test.Request("POST")
	resp := httptest.NewRecorder()
	w := harness.e.NewResponse(resp)

	// This event is really the only thing that separates "user not found"
	// from "bad password"
	var afterCalled bool
	harness.e.AuthEvents.After(engine.EventAuthFail, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
		afterCalled = true
		return false, nil
	})

	if err := harness.users.LoginPost(w, r); err != nil {
		t.Error(err)
	}

	if resp.Code != 200 {
		t.Error("wanted a 200:", resp.Code)
	}

	if harness.responder.Data[engine.DataErr] != "Invalid Credentials" {
		t.Error("wrong error:", harness.responder.Data)
	}

	if _, ok := harness.session.ClientValues[engine.SessionKey]; ok {
		t.Error("user should not be logged in")
	}

	if afterCalled {
		t.Error("after should not have been called")
	}
}

type testRedirector struct {
	Opts engine.RedirectOptions
}

func (r *testRedirector) Redirect(w http.ResponseWriter, req *http.Request, ro engine.RedirectOptions) error {
	r.Opts = ro
	if len(ro.RedirectPath) == 0 {
		panic("no redirect path on redirect call")
	}
	http.Redirect(w, req, ro.RedirectPath, ro.Code)
	return nil
}

type mockLogger struct{}

func (m mockLogger) Info(s string)  {}
func (m mockLogger) Error(s string) {}

func TestEngineMiddleware(t *testing.T) {
	t.Parallel()

	e := engine.New()
	e.Core.Logger = mockLogger{}
	e.Core.Server = &test.ServerStorer{
		Users: map[string]*test.User{
			"test@test.com": {},
		},
	}
	u := NewController(e)

	setupMore := func(mountPathed bool, requirements MWRequirements, failResponse MWRespondOnFailure) (*httptest.ResponseRecorder, bool, bool) {
		r := httptest.NewRequest("GET", "/super/secret", nil)
		rec := httptest.NewRecorder()
		w := e.NewResponse(rec)

		var err error
		r, err = e.LoadClientState(w, r)
		if err != nil {
			t.Fatal(err)
		}

		var mid func(http.Handler) http.Handler
		if !mountPathed {
			mid = u.AuthenticatedMiddleware(requirements, failResponse)
		} else {
			mid = u.AuthenticatedMountedMiddleware(true, requirements, failResponse)
		}
		var called, hadUser bool
		server := mid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			hadUser = r.Context().Value(CTXKeyUser) != nil
			w.WriteHeader(http.StatusOK)
		}))

		server.ServeHTTP(w, r)

		return rec, called, hadUser
	}

	t.Run("Accept", func(t *testing.T) {
		e.Core.SessionState = &test.ClientStateRW{
			ClientValues: map[string]string{engine.SessionKey: "test@test.com"},
		}

		_, called, hadUser := setupMore(false, RequireNone, RespondNotFound)

		if !called {
			t.Error("should have been called")
		}
		if !hadUser {
			t.Error("should have had user")
		}
	})
	t.Run("AcceptHalfAuth", func(t *testing.T) {
		e.Core.SessionState = &test.ClientStateRW{
			ClientValues: map[string]string{engine.SessionKey: "test@test.com", engine.SessionHalfAuthKey: "true"},
		}

		_, called, hadUser := setupMore(false, RequireNone, RespondNotFound)

		if !called {
			t.Error("should have been called")
		}
		if !hadUser {
			t.Error("should have had user")
		}
	})
	t.Run("RejectNotFound", func(t *testing.T) {
		e.Core.SessionState = test.NewClientRW()

		rec, called, hadUser := setupMore(false, RequireNone, RespondNotFound)

		if rec.Code != http.StatusNotFound {
			t.Error("wrong code:", rec.Code)
		}
		if called {
			t.Error("should not have been called")
		}
		if hadUser {
			t.Error("should not have had user")
		}
	})
	t.Run("RejectUnauthorized", func(t *testing.T) {
		e.Core.SessionState = test.NewClientRW()

		r := httptest.NewRequest("GET", "/super/secret", nil)
		rec := httptest.NewRecorder()
		w := e.NewResponse(rec)

		var err error
		r, err = e.LoadClientState(w, r)
		if err != nil {
			t.Fatal(err)
		}

		var mid func(http.Handler) http.Handler
		mid = u.AuthenticatedMiddleware(RequireNone, RespondUnauthorized)
		var called, hadUser bool
		server := mid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			hadUser = r.Context().Value(CTXKeyUser) != nil
			w.WriteHeader(http.StatusOK)
		}))

		server.ServeHTTP(w, r)

		if rec.Code != http.StatusUnauthorized {
			t.Error("wrong code:", rec.Code)
		}
		if called {
			t.Error("should not have been called")
		}
		if hadUser {
			t.Error("should not have had user")
		}
	})
	t.Run("RejectRedirect", func(t *testing.T) {
		redir := &testRedirector{}
		e.Core.Redirector = redir

		e.Core.SessionState = test.NewClientRW()

		_, called, hadUser := setupMore(false, RequireNone, RespondRedirect)

		if redir.Opts.Code != http.StatusTemporaryRedirect {
			t.Error("code was wrong:", redir.Opts.Code)
		}
		if redir.Opts.RedirectPath != "/login?redir=%2Fsuper%2Fsecret" {
			t.Error("redirect path was wrong:", redir.Opts.RedirectPath)
		}
		if called {
			t.Error("should not have been called")
		}
		if hadUser {
			t.Error("should not have had user")
		}
	})
	t.Run("RejectMountpathedRedirect", func(t *testing.T) {
		redir := &testRedirector{}
		e.Core.Redirector = redir

		e.Core.SessionState = test.NewClientRW()

		_, called, hadUser := setupMore(true, RequireNone, RespondRedirect)

		if redir.Opts.Code != http.StatusTemporaryRedirect {
			t.Error("code was wrong:", redir.Opts.Code)
		}
		if redir.Opts.RedirectPath != "/login?redir=%2Fsuper%2Fsecret" {
			t.Error("redirect path was wrong:", redir.Opts.RedirectPath)
		}
		if called {
			t.Error("should not have been called")
		}
		if hadUser {
			t.Error("should not have had user")
		}
	})
	t.Run("RequireFullAuth", func(t *testing.T) {
		e.Core.SessionState = &test.ClientStateRW{
			ClientValues: map[string]string{engine.SessionKey: "test@test.com", engine.SessionHalfAuthKey: "true"},
		}

		rec, called, hadUser := setupMore(false, RequireFullAuth, RespondNotFound)

		if rec.Code != http.StatusNotFound {
			t.Error("wrong code:", rec.Code)
		}
		if called {
			t.Error("should not have been called")
		}
		if hadUser {
			t.Error("should not have had user")
		}
	})
}
