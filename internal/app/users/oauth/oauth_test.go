package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func init() {
	exchanger = func(_ *oauth2.Config, _ context.Context, _ string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
		return testToken, nil
	}
}

var testProviders = map[string]engine.OAuth2Provider{
	"google": {
		OAuth2Config: &oauth2.Config{
			ClientID:     `jazz`,
			ClientSecret: `hands`,
			Scopes:       []string{`profile`, `email`},
			Endpoint:     google.Endpoint,
			// This is typically set by Init() but some tests rely on it's existence
			RedirectURL: "https://www.example.com/auth/oauth2/callback/google",
		},
		FindUserDetails:  GoogleUserDetails,
		AdditionalParams: url.Values{"include_requested_scopes": []string{"true"}},
	},
	"facebook": {
		OAuth2Config: &oauth2.Config{
			ClientID:     `jazz`,
			ClientSecret: `hands`,
			Scopes:       []string{`email`},
			Endpoint:     facebook.Endpoint,
			// This is typically set by Init() but some tests rely on it's existence
			RedirectURL: "https://www.example.com/auth/oauth2/callback/facebook",
		},
		FindUserDetails: FacebookUserDetails,
	},
}

var testToken = &oauth2.Token{
	AccessToken:  "token",
	TokenType:    "Bearer",
	RefreshToken: "refresh",
	Expiry:       time.Now().AddDate(0, 0, 1),
}

func TestInit(t *testing.T) {
	// No t.Parallel() since the cfg.RedirectURL is set in Init()

	e := engine.New()
	oauth := &OAuth2{}

	router := &test.Router{}
	e.Config.Modules.OAuth2Providers = testProviders
	e.Config.Core.Router = router
	e.Config.Core.ErrorHandler = &test.ErrorHandler{}

	e.Config.Paths.Mount = "/auth"
	e.Config.Paths.RootURL = "https://www.example.com"

	if err := oauth.Init(e); err != nil {
		t.Fatal(err)
	}

	gets := []string{
		"/oauth2/facebook", "/oauth2/callback/facebook",
		"/oauth2/google", "/oauth2/callback/google",
	}
	if err := router.HasGets(gets...); err != nil {
		t.Error(err)
	}
}

type testHarness struct {
	oauth *OAuth2
	e    *engine.Engine

	redirector *test.Redirector
	session    *test.ClientStateRW
	storer     *test.ServerStorer
}

func testSetup() *testHarness {
	harness := &testHarness{}

	harness.e = engine.New()
	harness.redirector = &test.Redirector{}
	harness.session = test.NewClientRW()
	harness.storer = test.NewServerStorer()

	harness.e.Modules.OAuth2Providers = testProviders

	harness.e.Config.Core.Logger = test.Logger{}
	harness.e.Config.Core.Redirector = harness.redirector
	harness.e.Config.Storage.SessionState = harness.session
	harness.e.Config.Storage.Server = harness.storer

	harness.oauth = &OAuth2{harness.e}

	return harness
}

func TestStart(t *testing.T) {
	t.Parallel()

	h := testSetup()

	rec := httptest.NewRecorder()
	w := h.e.NewResponse(rec)
	r := httptest.NewRequest("GET", "/oauth2/google?cake=yes&death=no", nil)

	if err := h.oauth.Start(w, r); err != nil {
		t.Error(err)
	}

	if h.redirector.Options.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", h.redirector.Options.Code)
	}

	redirectPathURL, err := url.Parse(h.redirector.Options.RedirectPath)
	if err != nil {
		t.Fatal(err)
	}
	query := redirectPathURL.Query()
	if state := query.Get("state"); len(state) == 0 {
		t.Error("our nonce should have been here")
	}
	if callback := query.Get("redirect_uri"); callback != "https://www.example.com/auth/oauth2/callback/google" {
		t.Error("callback was wrong:", callback)
	}
	if clientID := query.Get("client_id"); clientID != "jazz" {
		t.Error("clientID was wrong:", clientID)
	}
	if redirectPathURL.Host != "accounts.google.com" {
		t.Error("host was wrong:", redirectPathURL.Host)
	}

	if h.session.ClientValues[engine.SessionOAuth2State] != query.Get("state") {
		t.Error("the state should have been saved in the session")
	}
	if v := h.session.ClientValues[engine.SessionOAuth2Params]; v != `{"cake":"yes","death":"no"}` {
		t.Error("oauth2 session params are wrong:", v)
	}
}

func TestStartBadProvider(t *testing.T) {
	t.Parallel()

	h := testSetup()

	rec := httptest.NewRecorder()
	w := h.e.NewResponse(rec)
	r := httptest.NewRequest("GET", "/oauth2/test", nil)

	err := h.oauth.Start(w, r)
	if e := err.Error(); !strings.Contains(e, `provider "test" not found`) {
		t.Error("it should have errored:", e)
	}
}

func TestEnd(t *testing.T) {
	t.Parallel()

	h := testSetup()

	rec := httptest.NewRecorder()
	w := h.e.NewResponse(rec)

	h.session.ClientValues[engine.SessionOAuth2State] = "state"
	r, err := h.e.LoadClientState(w, httptest.NewRequest("GET", "/oauth2/callback/google?state=state", nil))
	if err != nil {
		t.Fatal(err)
	}

	if err := h.oauth.End(w, r); err != nil {
		t.Error(err)
	}

	w.WriteHeader(http.StatusOK) // Flush headers

	opts := h.redirector.Options
	if opts.Code != http.StatusTemporaryRedirect {
		t.Error("it should have redirected")
	}
	if opts.RedirectPath != "/" {
		t.Error("redir path was wrong:", opts.RedirectPath)
	}
	if s := h.session.ClientValues[engine.SessionKey]; s != "oauth2;;google;;id" {
		t.Error("session id should have been set:", s)
	}
}

func TestEndBadProvider(t *testing.T) {
	t.Parallel()

	h := testSetup()

	rec := httptest.NewRecorder()
	w := h.e.NewResponse(rec)
	r := httptest.NewRequest("GET", "/oauth2/callback/test", nil)

	err := h.oauth.End(w, r)
	if e := err.Error(); !strings.Contains(e, `provider "test" not found`) {
		t.Error("it should have errored:", e)
	}
}

func TestEndBadState(t *testing.T) {
	t.Parallel()

	h := testSetup()

	rec := httptest.NewRecorder()
	w := h.e.NewResponse(rec)
	r := httptest.NewRequest("GET", "/oauth2/callback/google", nil)

	err := h.oauth.End(w, r)
	if e := err.Error(); !strings.Contains(e, `oauth2 endpoint hit without session state`) {
		t.Error("it should have errored:", e)
	}

	h.session.ClientValues[engine.SessionOAuth2State] = "state"
	r, err = h.e.LoadClientState(w, httptest.NewRequest("GET", "/oauth2/callback/google?state=x", nil))
	if err != nil {
		t.Fatal(err)
	}
	if err := h.oauth.End(w, r); err != errOAuthStateValidation {
		t.Error("error was wrong:", err)
	}
}

func TestEndErrors(t *testing.T) {
	t.Parallel()

	h := testSetup()

	rec := httptest.NewRecorder()
	w := h.e.NewResponse(rec)

	h.session.ClientValues[engine.SessionOAuth2State] = "state"
	r, err := h.e.LoadClientState(w, httptest.NewRequest("GET", "/oauth2/callback/google?state=state&error=badtimes&error_reason=reason", nil))
	if err != nil {
		t.Fatal(err)
	}

	if err := h.oauth.End(w, r); err != nil {
		t.Error(err)
	}

	opts := h.redirector.Options
	if opts.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", opts.Code)
	}
	if opts.RedirectPath != "/" {
		t.Error("path was wrong:", opts.RedirectPath)
	}
}

func TestEndHandling(t *testing.T) {
	t.Parallel()

	t.Run("AfterOAuth2Fail", func(t *testing.T) {
		h := testSetup()

		rec := httptest.NewRecorder()
		w := h.e.NewResponse(rec)

		h.session.ClientValues[engine.SessionOAuth2State] = "state"
		r, err := h.e.LoadClientState(w, httptest.NewRequest("GET", "/oauth2/callback/google?state=state&error=badtimes&error_reason=reason", nil))
		if err != nil {
			t.Fatal(err)
		}

		called := false
		h.e.Events.After(engine.EventOAuth2Fail, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			called = true
			return true, nil
		})

		if err := h.oauth.End(w, r); err != nil {
			t.Error(err)
		}

		if !called {
			t.Error("it should have been called")
		}
		if h.redirector.Options.Code != 0 {
			t.Error("it should not have tried to redirect")
		}
	})
	t.Run("BeforeOAuth2", func(t *testing.T) {
		h := testSetup()

		rec := httptest.NewRecorder()
		w := h.e.NewResponse(rec)

		h.session.ClientValues[engine.SessionOAuth2State] = "state"
		r, err := h.e.LoadClientState(w, httptest.NewRequest("GET", "/oauth2/callback/google?state=state", nil))
		if err != nil {
			t.Fatal(err)
		}

		called := false
		h.e.Events.Before(engine.EventOAuth2, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			called = true
			return true, nil
		})

		if err := h.oauth.End(w, r); err != nil {
			t.Error(err)
		}

		w.WriteHeader(http.StatusOK) // Flush headers

		if !called {
			t.Error("it should have been called")
		}
		if h.redirector.Options.Code != 0 {
			t.Error("it should not have tried to redirect")
		}
		if len(h.session.ClientValues[engine.SessionKey]) != 0 {
			t.Error("should have not logged the user in")
		}
	})

	t.Run("AfterOAuth2", func(t *testing.T) {
		h := testSetup()

		rec := httptest.NewRecorder()
		w := h.e.NewResponse(rec)

		h.session.ClientValues[engine.SessionOAuth2State] = "state"
		r, err := h.e.LoadClientState(w, httptest.NewRequest("GET", "/oauth2/callback/google?state=state", nil))
		if err != nil {
			t.Fatal(err)
		}

		called := false
		h.e.Events.After(engine.EventOAuth2, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			called = true
			return true, nil
		})

		if err := h.oauth.End(w, r); err != nil {
			t.Error(err)
		}

		w.WriteHeader(http.StatusOK) // Flush headers

		if !called {
			t.Error("it should have been called")
		}
		if h.redirector.Options.Code != 0 {
			t.Error("it should not have tried to redirect")
		}
		if s := h.session.ClientValues[engine.SessionKey]; s != "oauth2;;google;;id" {
			t.Error("session id should have been set:", s)
		}
	})
}
