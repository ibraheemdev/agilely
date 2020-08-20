package logoutable

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
	"github.com/ibraheemdev/agilely/test/authboss"
)

func TestLogout(t *testing.T) {
	t.Parallel()

	ab := authboss.New()

	router := &authboss_test.Router{}
	errHandler := &authboss_test.ErrorHandler{}
	ab.Config.Core.Router = router
	ab.Config.Core.ErrorHandler = errHandler

	l := &Logout{}
	if err := l.Init(ab); err != nil {
		t.Fatal(err)
	}

	if err := router.HasDeletes("/logout"); err != nil {
		t.Error(err)
	}
}

func TestLogoutRoutes(t *testing.T) {
	t.Parallel()

	ab := authboss.New()
	router := &authboss_test.Router{}
	errHandler := &authboss_test.ErrorHandler{}
	ab.Config.Core.Router = router
	ab.Config.Core.ErrorHandler = errHandler

	l := &Logout{}

	if err := l.Init(ab); err != nil {
		t.Error("should have failed to register the route")
	}
	if err := router.HasDeletes("/logout"); err != nil {
		t.Error(err)
	}
}

type testHarness struct {
	logout *Logout
	ab     *authboss.Authboss

	redirector *authboss_test.Redirector
	session    *authboss_test.ClientStateRW
	cookies    *authboss_test.ClientStateRW
	storer     *authboss_test.ServerStorer
}

func testSetup() *testHarness {
	harness := &testHarness{}

	harness.ab = authboss.New()
	harness.redirector = &authboss_test.Redirector{}
	harness.session = authboss_test.NewClientRW()
	harness.cookies = authboss_test.NewClientRW()
	harness.storer = authboss_test.NewServerStorer()

	harness.ab.Paths.LogoutOK = "/logout/ok"

	harness.ab.Config.Core.Logger = authboss_test.Logger{}
	harness.ab.Config.Core.Redirector = harness.redirector
	harness.ab.Config.Storage.SessionState = harness.session
	harness.ab.Config.Storage.CookieState = harness.cookies
	harness.ab.Config.Storage.Server = harness.storer

	harness.logout = &Logout{harness.ab}

	return harness
}

func TestLogoutLogout(t *testing.T) {
	t.Parallel()

	h := testSetup()

	h.session.ClientValues[authboss.SessionKey] = "test@test.com"
	h.session.ClientValues[authboss.SessionHalfAuthKey] = "true"
	h.session.ClientValues[authboss.SessionLastAction] = time.Now().UTC().Format(time.RFC3339)
	h.cookies.ClientValues[authboss.CookieRemember] = "token"

	r := authboss_test.Request("POST")
	resp := httptest.NewRecorder()
	w := h.ab.NewResponse(resp)

	// This enables the logging portion
	// which is debatable-y not useful in a log out method
	user := &authboss_test.User{Email: "test@test.com"}
	r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyUser, user))

	var err error
	r, err = h.ab.LoadClientState(w, r)
	if err != nil {
		t.Error(err)
	}

	if err := h.logout.Logout(w, r); err != nil {
		t.Fatal(err)
	}

	if resp.Code != http.StatusTemporaryRedirect {
		t.Error("response code wrong:", resp.Code)
	}
	if h.redirector.Options.RedirectPath != "/logout/ok" {
		t.Error("redirect path was wrong:", h.redirector.Options.RedirectPath)
	}

	if _, ok := h.session.ClientValues[authboss.SessionKey]; ok {
		t.Error("want session key gone")
	}
	if _, ok := h.session.ClientValues[authboss.SessionHalfAuthKey]; ok {
		t.Error("want session half auth key gone")
	}
	if _, ok := h.session.ClientValues[authboss.SessionLastAction]; ok {
		t.Error("want session last action")
	}
	if _, ok := h.cookies.ClientValues[authboss.CookieRemember]; ok {
		t.Error("want remember me cookies gone")
	}
}
