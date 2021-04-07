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

type testLogoutHarness struct {
	users *Users
	e     *engine.Engine

	redirector *test.Redirector
	session    *test.ClientStateRW
	cookies    *test.ClientStateRW
	storer     *test.ServerStorer
}

func testLogoutSetup() *testLogoutHarness {
	harness := &testLogoutHarness{}

	harness.e = engine.New()
	harness.redirector = &test.Redirector{}
	harness.session = test.NewClientRW()
	harness.cookies = test.NewClientRW()
	harness.storer = test.NewServerStorer()

	harness.e.Core.Logger = test.Logger{}
	harness.e.Core.Redirector = harness.redirector
	harness.e.Core.SessionState = harness.session
	harness.e.Core.CookieState = harness.cookies
	harness.e.Core.Database = harness.storer

	harness.users = NewController(harness.e)

	return harness
}

func TestLogoutLogout(t *testing.T) {
	t.Parallel()

	h := testLogoutSetup()

	h.session.ClientValues[engine.SessionKey] = "test@test.com"
	h.session.ClientValues[SessionHalfAuthKey] = "true"
	h.session.ClientValues[engine.SessionLastAction] = time.Now().UTC().Format(time.RFC3339)
	h.cookies.ClientValues[CookieRemember] = "token"

	r := test.Request("POST")
	resp := httptest.NewRecorder()
	w := h.e.NewResponse(resp)

	// This enables the logging portion
	// which is debatable-y not useful in a log out method
	user := &test.User{Email: "test@test.com"}
	r = r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))

	var err error
	r, err = h.e.LoadClientState(w, r)
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
	if _, ok := h.session.ClientValues[SessionHalfAuthKey]; ok {
		t.Error("want session half auth key gone")
	}
	if _, ok := h.session.ClientValues[engine.SessionLastAction]; ok {
		t.Error("want session last action")
	}
	if _, ok := h.cookies.ClientValues[CookieRemember]; ok {
		t.Error("want remember me cookies gone")
	}
}
