package users

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
)

func TestRegisterInit(t *testing.T) {
	t.Parallel()

	e := engine.New()

	router := &test.Router{}
	renderer := &test.Renderer{}
	errHandler := &test.ErrorHandler{}
	e.Config.Core.Router = router
	e.Config.Core.ViewRenderer = renderer
	e.Config.Core.ErrorHandler = errHandler
	e.Config.Storage.Server = &test.ServerStorer{}

	u := &Users{}
	if err := u.InitRegister(); err != nil {
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

	ab := engine.New()
	responder := &test.Responder{}
	ab.Config.Core.Responder = responder

	u := &Users{ab}
	if err := u.GetRegister(nil, nil); err != nil {
		t.Error(err)
	}

	if responder.Page != PageRegister {
		t.Error("wanted login page, got:", responder.Page)
	}

	if responder.Status != http.StatusOK {
		t.Error("wanted ok status, got:", responder.Status)
	}
}

type testRegisterHarness struct {
	users *Users
	ab    *engine.Engine

	bodyReader *test.BodyReader
	responder  *test.Responder
	redirector *test.Redirector
	session    *test.ClientStateRW
	storer     *test.ServerStorer
}

func testRegisterSetup() *testRegisterHarness {
	harness := &testRegisterHarness{}

	harness.ab = engine.New()
	harness.bodyReader = &test.BodyReader{}
	harness.redirector = &test.Redirector{}
	harness.responder = &test.Responder{}
	harness.session = test.NewClientRW()
	harness.storer = test.NewServerStorer()

	harness.ab.Config.Core.BodyReader = harness.bodyReader
	harness.ab.Config.Core.Logger = test.Logger{}
	harness.ab.Config.Core.Responder = harness.responder
	harness.ab.Config.Core.Redirector = harness.redirector
	harness.ab.Config.Storage.SessionState = harness.session
	harness.ab.Config.Storage.Server = harness.storer

	harness.users = &Users{harness.ab}

	return harness
}

func TestRegisterPostSuccess(t *testing.T) {
	t.Parallel()

	setupMore := func(harness *testRegisterHarness) *testRegisterHarness {
		harness.ab.Modules.RegisterPreserveFields = []string{"email", "another"}
		harness.bodyReader.Return = test.ArbValues{
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
		h := setupMore(testRegisterSetup())

		r := test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.ab.NewResponse(resp)

		if err := h.users.PostRegister(w, r); err != nil {
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

		if h.session.ClientValues[engine.SessionKey] != "test@test.com" {
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
		h := setupMore(testRegisterSetup())

		var afterCalled bool
		h.ab.Events.After(engine.EventRegister, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			w.WriteHeader(http.StatusTeapot)
			afterCalled = true
			return true, nil
		})

		r := test.Request("POST")
		resp := httptest.NewRecorder()
		w := h.ab.NewResponse(resp)

		if err := h.users.PostRegister(w, r); err != nil {
			t.Error(err)
		}

		user, ok := h.storer.Users["test@test.com"]
		if !ok {
			t.Error("user was not persisted in the DB")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("hello world")); err != nil {
			t.Error("password was not properly encrypted:", err)
		}

		if val, ok := h.session.ClientValues[engine.SessionKey]; ok {
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

	h := testRegisterSetup()

	// Ensure the below is sorted, the sort normally happens in Init()
	// that we don't call
	h.ab.Modules.RegisterPreserveFields = []string{"another", "email"}
	h.bodyReader.Return = test.ArbValues{
		Values: map[string]string{
			"email":    "test@test.com",
			"password": "hello world",
			"another":  "value",
		},
		Errors: []error{
			errors.New("bad password"),
		},
	}

	r := test.Request("POST")
	resp := httptest.NewRecorder()
	w := h.ab.NewResponse(resp)

	if err := h.users.PostRegister(w, r); err != nil {
		t.Error(err)
	}

	if h.responder.Status != http.StatusOK {
		t.Error("wrong status:", h.responder.Status)
	}
	if h.responder.Page != PageRegister {
		t.Error("rendered wrong page:", h.responder.Page)
	}

	errList := h.responder.Data[engine.DataValidation].(map[string][]string)
	if e := errList[""][0]; e != "bad password" {
		t.Error("validation error wrong:", e)
	}

	intfD, ok := h.responder.Data[engine.DataPreserve]
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

	h := testRegisterSetup()

	// Ensure the below is sorted, the sort normally happens in Init()
	// that we don't call
	h.ab.Modules.RegisterPreserveFields = []string{"another", "email"}
	h.storer.Users["test@test.com"] = &test.User{}
	h.bodyReader.Return = test.ArbValues{
		Values: map[string]string{
			"email":    "test@test.com",
			"password": "hello world",
			"another":  "value",
		},
	}

	r := test.Request("POST")
	resp := httptest.NewRecorder()
	w := h.ab.NewResponse(resp)

	if err := h.users.PostRegister(w, r); err != nil {
		t.Error(err)
	}

	if h.responder.Status != http.StatusOK {
		t.Error("wrong status:", h.responder.Status)
	}
	if h.responder.Page != PageRegister {
		t.Error("rendered wrong page:", h.responder.Page)
	}

	errList := h.responder.Data[engine.DataValidation].(map[string][]string)
	if e := errList[""][0]; e != "user already exists" {
		t.Error("validation error wrong:", e)
	}

	intfD, ok := h.responder.Data[engine.DataPreserve]
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
