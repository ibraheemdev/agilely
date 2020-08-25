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

func TestLogout(t *testing.T) {
	t.Parallel()

	ab := engine.New()

	router := &test.Router{}
	errHandler := &test.ErrorHandler{}
	ab.Config.Core.Router = router
	ab.Config.Core.ErrorHandler = errHandler

	u := &Users{}
	if err := u.InitLogout(); err != nil {
		t.Fatal(err)
	}

	if err := router.HasDeletes("/logout"); err != nil {
		t.Error(err)
	}
}

func TestLogoutRoutes(t *testing.T) {
	t.Parallel()

	ab := engine.New()
	router := &test.Router{}
	errHandler := &test.ErrorHandler{}
	ab.Config.Core.Router = router
	ab.Config.Core.ErrorHandler = errHandler

	u := &Users{}

	if err := u.InitLogout(); err != nil {
		t.Error("should have failed to register the route")
	}
	if err := router.HasDeletes("/logout"); err != nil {
		t.Error(err)
	}
}

type testLogoutHarness struct {
	users *Users
	ab    *engine.Engine

	redirector *test.Redirector
	session    *test.ClientStateRW
	cookies    *test.ClientStateRW
	storer     *test.ServerStorer
}

func testLogoutSetup() *testLogoutHarness {
	harness := &testLogoutHarness{}

	harness.ab = engine.New()
	harness.redirector = &test.Redirector{}
	harness.session = test.NewClientRW()
	harness.cookies = test.NewClientRW()
	harness.storer = test.NewServerStorer()

	harness.ab.Config.Core.Logger = test.Logger{}
	harness.ab.Config.Core.Redirector = harness.redirector
	harness.ab.Config.Storage.SessionState = harness.session
	harness.ab.Config.Storage.CookieState = harness.cookies
	harness.ab.Config.Storage.Server = harness.storer

	harness.users = &Users{harness.ab}

	return harness
}

func TestLogoutLogout(t *testing.T) {
	t.Parallel()

	h := testLogoutSetup()

	h.session.ClientValues[engine.SessionKey] = "test@test.com"
	h.session.ClientValues[engine.SessionHalfAuthKey] = "true"
	h.session.ClientValues[engine.SessionLastAction] = time.Now().UTC().Format(time.RFC3339)
	h.cookies.ClientValues[engine.CookieRemember] = "token"

	r := test.Request("POST")
	resp := httptest.NewRecorder()
	w := h.ab.NewResponse(resp)

	// This enables the logging portion
	// which is debatable-y not useful in a log out method
	user := &test.User{Email: "test@test.com"}
	r = r.WithContext(context.WithValue(r.Context(), engine.CTXKeyUser, user))

	var err error
	r, err = h.ab.LoadClientState(w, r)
	if err != nil {
		t.Error(err)
	}

	if err := h.users.Logout(w, r); err != nil {
		t.Fatal(err)
	}

	if resp.Code != http.StatusTemporaryRedirect {
		t.Error("response code wrong:", resp.Code)
	}
	if h.redirector.Options.RedirectPath != "/" {
		t.Error("redirect path was wrong:", h.redirector.Options.RedirectPath)
	}

	if _, ok := h.session.ClientValues[engine.SessionKey]; ok {
		t.Error("want session key gone")
	}
	if _, ok := h.session.ClientValues[engine.SessionHalfAuthKey]; ok {
		t.Error("want session half auth key gone")
	}
	if _, ok := h.session.ClientValues[engine.SessionLastAction]; ok {
		t.Error("want session last action")
	}
	if _, ok := h.cookies.ClientValues[engine.CookieRemember]; ok {
		t.Error("want remember me cookies gone")
	}
}
