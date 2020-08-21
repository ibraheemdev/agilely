package registerable

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
	"github.com/ibraheemdev/agilely/test/authboss"
)

func TestRegisterInit(t *testing.T) {
	t.Parallel()

	ab := authboss.New()

	router := &authboss_test.Router{}
	renderer := &authboss_test.Renderer{}
	errHandler := &authboss_test.ErrorHandler{}
	ab.Config.Core.Router = router
	ab.Config.Core.ViewRenderer = renderer
	ab.Config.Core.ErrorHandler = errHandler
	ab.Config.Storage.Server = &authboss_test.ServerStorer{}

	reg := &Register{}
	if err := reg.Init(ab); err != nil {
		t.Fatal(err)
	}

	if err := renderer.HasLoadedViews(PageRegister); err != nil {
		t.Error(err)
	}

	if err := router.HasGets("/register"); err != nil {
		t.Error(err)
	}
	if err := router.HasPosts("/register"); err != nil {
		t.Error(err)
	}
}

func TestRegisterGet(t *testing.T) {
	t.Parallel()

	ab := authboss.New()
	responder := &authboss_test.Responder{}
	ab.Config.Core.Responder = responder

	a := &Register{ab}
	if err := a.Get(nil, nil); err != nil {
		t.Error(err)
	}

	if responder.Page != PageRegister {
		t.Error("wanted login page, got:", responder.Page)
	}

	if responder.Status != http.StatusOK {
		t.Error("wanted ok status, got:", responder.Status)
	}
}

type testHarness struct {
	reg *Register
	ab  *authboss.Authboss

	bodyReader *authboss_test.BodyReader
	responder  *authboss_test.Responder
	redirector *authboss_test.Redirector
	session    *authboss_test.ClientStateRW
	storer     *authboss_test.ServerStorer
}

func testSetup() *testHarness {
	harness := &testHarness{}

	harness.ab = authboss.New()
	harness.bodyReader = &authboss_test.BodyReader{}
	harness.redirector = &authboss_test.Redirector{}
	harness.responder = &authboss_test.Responder{}
	harness.session = authboss_test.NewClientRW()
	harness.storer = authboss_test.NewServerStorer()

	harness.ab.Config.Core.BodyReader = harness.bodyReader
	harness.ab.Config.Core.Logger = authboss_test.Logger{}
	harness.ab.Config.Core.Responder = harness.responder
	harness.ab.Config.Core.Redirector = harness.redirector
	harness.ab.Config.Storage.SessionState = harness.session
	harness.ab.Config.Storage.Server = harness.storer

	harness.reg = &Register{harness.ab}

	return harness
}

func TestRegisterPostSuccess(t *testing.T) {
	t.Parallel()

	setupMore := func(harness *testHarness) *testHarness {
		harness.ab.Modules.RegisterPreserveFields = []string{"email", "another"}
		harness.bodyReader.Return = authboss_test.ArbValues{
			Values: map[string]string{
				"email":    "test@test.com",
				"password": "hello world",
				"another":  "value",
			},
		}

		return harness
	}

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		h := setupMore(testSetup())

		r := authboss_test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.ab.NewResponse(resp)

		if err := h.reg.Post(w, r); err != nil {
			t.Error(err)
		}

		user, ok := h.storer.Users["test@test.com"]
		if !ok {
			t.Error("user was not persisted in the DB")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("hello world")); err != nil {
			t.Error("password was not properly encrypted:", err)
		}

		if user.Arbitrary["another"] != "value" {
			t.Error("arbitrary values not saved")
		}

		if h.session.ClientValues[authboss.SessionKey] != "test@test.com" {
			t.Error("user should have been logged in:", h.session.ClientValues)
		}

		if resp.Code != http.StatusTemporaryRedirect {
			t.Error("code was wrong:", resp.Code)
		}
		if h.redirector.Options.RedirectPath != "/login" {
			t.Error("redirect path was wrong:", h.redirector.Options.RedirectPath)
		}
	})

	t.Run("handledAfter", func(t *testing.T) {
		t.Parallel()
		h := setupMore(testSetup())

		var afterCalled bool
		h.ab.Events.After(authboss.EventRegister, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			w.WriteHeader(http.StatusTeapot)
			afterCalled = true
			return true, nil
		})

		r := authboss_test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.ab.NewResponse(resp)

		if err := h.reg.Post(w, r); err != nil {
			t.Error(err)
		}

		user, ok := h.storer.Users["test@test.com"]
		if !ok {
			t.Error("user was not persisted in the DB")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("hello world")); err != nil {
			t.Error("password was not properly encrypted:", err)
		}

		if val, ok := h.session.ClientValues[authboss.SessionKey]; ok {
			t.Error("user should not have been logged in:", val)
		}

		if resp.Code != http.StatusTeapot {
			t.Error("code was wrong:", resp.Code)
		}

		if !afterCalled {
			t.Error("the after handler should have been called")
		}
	})
}

func TestRegisterPostValidationFailure(t *testing.T) {
	t.Parallel()

	h := testSetup()

	// Ensure the below is sorted, the sort normally happens in Init()
	// that we don't call
	h.ab.Modules.RegisterPreserveFields = []string{"another", "email"}
	h.bodyReader.Return = authboss_test.ArbValues{
		Values: map[string]string{
			"email":    "test@test.com",
			"password": "hello world",
			"another":  "value",
		},
		Errors: []error{
			errors.New("bad password"),
		},
	}

	r := authboss_test.Request("POST")
	resp := httptest.NewRecorder()
	w := h.ab.NewResponse(resp)

	if err := h.reg.Post(w, r); err != nil {
		t.Error(err)
	}

	if h.responder.Status != http.StatusOK {
		t.Error("wrong status:", h.responder.Status)
	}
	if h.responder.Page != PageRegister {
		t.Error("rendered wrong page:", h.responder.Page)
	}

	errList := h.responder.Data[authboss.DataValidation].(map[string][]string)
	if e := errList[""][0]; e != "bad password" {
		t.Error("validation error wrong:", e)
	}

	intfD, ok := h.responder.Data[authboss.DataPreserve]
	if !ok {
		t.Fatal("there was no preserved data")
	}

	d := intfD.(map[string]string)
	if d["email"] != "test@test.com" {
		t.Error("e-mail was not preserved:", d)
	} else if d["another"] != "value" {
		t.Error("another value was not preserved", d)
	} else if _, ok = d["password"]; ok {
		t.Error("password was preserved", d)
	}
}

func TestRegisterPostUserExists(t *testing.T) {
	t.Parallel()

	h := testSetup()

	// Ensure the below is sorted, the sort normally happens in Init()
	// that we don't call
	h.ab.Modules.RegisterPreserveFields = []string{"another", "email"}
	h.storer.Users["test@test.com"] = &authboss_test.User{}
	h.bodyReader.Return = authboss_test.ArbValues{
		Values: map[string]string{
			"email":    "test@test.com",
			"password": "hello world",
			"another":  "value",
		},
	}

	r := authboss_test.Request("POST")
	resp := httptest.NewRecorder()
	w := h.ab.NewResponse(resp)

	if err := h.reg.Post(w, r); err != nil {
		t.Error(err)
	}

	if h.responder.Status != http.StatusOK {
		t.Error("wrong status:", h.responder.Status)
	}
	if h.responder.Page != PageRegister {
		t.Error("rendered wrong page:", h.responder.Page)
	}

	errList := h.responder.Data[authboss.DataValidation].(map[string][]string)
	if e := errList[""][0]; e != "user already exists" {
		t.Error("validation error wrong:", e)
	}

	intfD, ok := h.responder.Data[authboss.DataPreserve]
	if !ok {
		t.Fatal("there was no preserved data")
	}

	d := intfD.(map[string]string)
	if d["email"] != "test@test.com" {
		t.Error("e-mail was not preserved:", d)
	} else if d["another"] != "value" {
		t.Error("another value was not preserved", d)
	} else if _, ok = d["password"]; ok {
		t.Error("password was preserved", d)
	}
}

func TestHasString(t *testing.T) {
	t.Parallel()

	strs := []string{"b", "c", "d", "e"}

	if !hasString(strs, "b") {
		t.Error("should have a")
	}
	if !hasString(strs, "e") {
		t.Error("should have d")
	}

	if hasString(strs, "a") {
		t.Error("should not have a")
	}
	if hasString(strs, "f") {
		t.Error("should not have f")
	}
}
