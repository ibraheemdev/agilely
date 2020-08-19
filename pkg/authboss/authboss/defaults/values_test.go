package defaults

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
	"github.com/ibraheemdev/agilely/test"
)

func TestHTTPBodyReaderLogin(t *testing.T) {
	t.Parallel()

	h := NewHTTPBodyReader(false, false)
	r := test.Request("POST", "email", "john@john.john", "password", "flowers")

	validator, err := h.Read("login.html.tpl", r)
	if err != nil {
		t.Error(err)
	}

	uv := validator.(authboss.UserValuer)
	if "john@john.john" != uv.GetPID() {
		t.Error("wrong e-mail:", uv.GetPID())
	}
	if "flowers" != uv.GetPassword() {
		t.Error("wrong password:", uv.GetPassword())
	}
}

func TestHTTPBodyReaderJSON(t *testing.T) {
	t.Parallel()

	h := NewHTTPBodyReader(true, false)
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"email":"john@john.john","password":"flowers"}`))

	validator, err := h.Read("login.html.tpl", r)
	if err != nil {
		t.Error(err)
	}

	uv := validator.(authboss.UserValuer)
	if "john@john.john" != uv.GetPID() {
		t.Error("wrong e-mail:", uv.GetPID())
	}
	if "flowers" != uv.GetPassword() {
		t.Error("wrong password:", uv.GetPassword())
	}
}

func TestHTTPBodyReaderConfirm(t *testing.T) {
	t.Parallel()

	h := NewHTTPBodyReader(false, false)
	r := test.Request("POST", FormValueConfirm, "token")

	validator, err := h.Read("mailer/confirm.html.tpl", r)
	if err != nil {
		t.Error(err)
	}

	cv := validator.(authboss.ConfirmValuer)
	if "token" != cv.GetToken() {
		t.Error("token was wrong:", cv.GetToken())
	}
}

func TestHTTPBodyReaderRecoverStart(t *testing.T) {
	t.Parallel()

	h := NewHTTPBodyReader(false, false)
	r := test.Request("POST", FormValueEmail, "email")

	validator, err := h.Read("recover_start.html.tpl", r)
	if err != nil {
		t.Error(err)
	}

	rsv := validator.(authboss.RecoverStartValuer)
	if pid := rsv.GetPID(); pid != "email" {
		t.Error("token was wrong:", pid)
	}
}

func TestHTTPBodyReaderRecoverMiddle(t *testing.T) {
	t.Parallel()

	h := NewHTTPBodyReader(false, false)
	r := httptest.NewRequest("GET", "/?token=token", nil)

	validator, err := h.Read("recover_middle.html.tpl", r)
	if err != nil {
		t.Error(err)
	}

	rmv := validator.(authboss.RecoverMiddleValuer)
	if token := rmv.GetToken(); token != "token" {
		t.Error("token was wrong:", token)
	}
}

func TestHTTPBodyReaderRecoverEnd(t *testing.T) {
	t.Parallel()

	h := NewHTTPBodyReader(false, false)
	r := test.Request("POST", "token", "token", "password", "password")

	validator, err := h.Read("recover_end.html.tpl", r)
	if err != nil {
		t.Error(err)
	}

	rmv := validator.(authboss.RecoverEndValuer)
	if token := rmv.GetToken(); token != "token" {
		t.Error("token was wrong:", token)
	}
	if password := rmv.GetPassword(); password != "password" {
		t.Error("password was wrong:", password)
	}
}

func TestHTTPBodyReaderRegister(t *testing.T) {
	t.Parallel()

	h := NewHTTPBodyReader(false, false)
	h.Whitelist["register.html.tpl"] = []string{"address"}
	r := test.Request("POST", "email", "a@a.com", "password", "1234", "address", "555 go street")

	validator, err := h.Read("register.html.tpl", r)
	if err != nil {
		t.Error(err)
	}

	rv := validator.(authboss.UserValuer)
	if pid := rv.GetPID(); pid != "a@a.com" {
		t.Error("pid was wrong:", pid)
	}
	if password := rv.GetPassword(); password != "1234" {
		t.Error("password was wrong:", password)
	}

	arb := validator.(authboss.ArbitraryValuer)
	values := arb.GetValues()
	if address := values["address"]; address != "555 go street" {
		t.Errorf("address was wrong, got: %s", values)
	}
}
