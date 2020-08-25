package users

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/test"
)

const (
	testSelector = `rnaGE8TDilrINHPxq/2xNU1FUTzsUSX8FvN5YzooyyWKk88fw1DjjbKBRGFtGew9OeZ+xeCC4mslfvQQMYspIg==`
	testVerifier = `W1Mz30QhavVM4d8jKaFtxGBfb4GX+fOn7V0Pc1WeftgtyOtY5OX7sY9gIeY5CIY4n8LvfWy14W7/6rs2KO9pgA==`
	testToken    = `w5OZ51E61Q6wsJOVr9o7KmyepP7Od5VBHQ1ADDUBkiGGMjKfnMFPjtvNpLjLKJqffw72KWZzNLj0Cs8wqywdEQ==`
)

type testRecoverHarness struct {
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

func testRecoverSetup() *testRecoverHarness {
	harness := &testRecoverHarness{}

	harness.e = engine.New()
	harness.bodyReader = &test.BodyReader{}
	harness.mailer = &test.Emailer{}
	harness.redirector = &test.Redirector{}
	harness.renderer = &test.Renderer{}
	harness.responder = &test.Responder{}
	harness.session = test.NewClientRW()
	harness.storer = test.NewServerStorer()

	harness.e.Config.Core.BodyReader = harness.bodyReader
	harness.e.Config.Core.Logger = test.Logger{}
	harness.e.Config.Core.Mailer = harness.mailer
	harness.e.Config.Core.Redirector = harness.redirector
	harness.e.Config.Core.MailRenderer = harness.renderer
	harness.e.Config.Core.Responder = harness.responder
	harness.e.Config.Storage.SessionState = harness.session
	harness.e.Config.Storage.Server = harness.storer

	harness.users = NewController(harness.e)

	return harness
}

func TestStartGet(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	r := test.Request("GET")
	w := httptest.NewRecorder()

	if err := h.users.StartGetRecover(w, r); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusOK {
		t.Error("code was wrong:", w.Code)
	}
	if h.responder.Page != PageRecoverStart {
		t.Error("page was wrong:", h.responder.Page)
	}
	if h.responder.Data != nil {
		t.Error("expected no data:", h.responder.Data)
	}
}

func TestStartPostSuccess(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		PID: "test@test.com",
	}
	h.storer.Users["test@test.com"] = &test.User{
		Email:    "test@test.com",
		Password: "i can't recall, doesn't seem like something bcrypted though",
	}

	r := test.Request("GET")
	w := httptest.NewRecorder()

	if err := h.users.StartPostRecover(w, r); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", w.Code)
	}
	if h.redirector.Options.RedirectPath != "/" {
		t.Error("page was wrong:", h.responder.Page)
	}
	if len(h.redirector.Options.Success) == 0 {
		t.Error("expected a nice success message")
	}

	time.Sleep(time.Second)
	if h.mailer.Email.To[0] != "test@test.com" {
		t.Error("e-mail to address is wrong:", h.mailer.Email.To)
	}
	if !strings.HasSuffix(h.mailer.Email.Subject, "Password Reset") {
		t.Error("e-mail subject line is wrong:", h.mailer.Email.Subject)
	}
	if len(h.renderer.Data[DataRecoverURL].(string)) == 0 {
		t.Errorf("the renderer's url in data was missing: %#v", h.renderer.Data)
	}
}

func TestStartPostFailure(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		PID: "test@test.com",
	}

	r := test.Request("GET")
	w := httptest.NewRecorder()

	if err := h.users.StartPostRecover(w, r); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", w.Code)
	}
	if h.redirector.Options.RedirectPath != "/" {
		t.Error("page was wrong:", h.responder.Page)
	}
	if len(h.redirector.Options.Success) == 0 {
		t.Error("expected a nice success message")
	}

	if len(h.mailer.Email.To) != 0 {
		t.Error("should not have sent an e-mail out!")
	}
}

func TestEndGet(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		Token: "abcd",
	}

	r := test.Request("GET")
	w := httptest.NewRecorder()

	if err := h.users.EndGetRecover(w, r); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusOK {
		t.Error("code was wrong:", w.Code)
	}
	if h.responder.Page != PageRecoverEnd {
		t.Error("page was wrong:", h.responder.Page)
	}
	if h.responder.Data[DataRecoverToken].(string) != "abcd" {
		t.Errorf("recovery token is wrong: %#v", h.responder.Data)
	}
}

func TestEndPostSuccess(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		Token: testToken,
	}
	h.storer.Users["test@test.com"] = &test.User{
		Email:              "test@test.com",
		Password:           "to-overwrite",
		RecoverSelector:    testSelector,
		RecoverVerifier:    testVerifier,
		RecoverTokenExpiry: time.Now().UTC().AddDate(0, 0, 1),
	}

	r := test.Request("POST")
	w := httptest.NewRecorder()

	if err := h.users.EndPostRecover(w, r); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", w.Code)
	}
	if p := h.redirector.Options.RedirectPath; p != "/" {
		t.Error("path was wrong:", p)
	}
	if len(h.session.ClientValues[engine.SessionKey]) != 0 {
		t.Error("should not have logged in the user")
	}
	if !strings.Contains(h.redirector.Options.Success, "updated password") {
		t.Error("should talk about recovering the password")
	}
	if strings.Contains(h.redirector.Options.Success, "logged in") {
		t.Error("should not talk about logging in")
	}
}

func TestEndPostSuccessLogin(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		Token: testToken,
	}
	h.storer.Users["test@test.com"] = &test.User{
		Email:              "test@test.com",
		Password:           "to-overwrite",
		RecoverSelector:    testSelector,
		RecoverVerifier:    testVerifier,
		RecoverTokenExpiry: time.Now().UTC().AddDate(0, 0, 1),
	}

	r := test.Request("POST")
	w := httptest.NewRecorder()

	if err := h.users.EndPostRecover(h.e.NewResponse(w), r); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusTemporaryRedirect {
		t.Error("code was wrong:", w.Code)
	}
	if p := h.redirector.Options.RedirectPath; p != "/" {
		t.Error("path was wrong:", p)
	}
	if len(h.session.ClientValues[engine.SessionKey]) != 0 {
		t.Error("it should not have logged in the user")
	}
	if h.redirector.Options.Success != "Successfully updated password" {
		t.Errorf("wrong success message: %s", h.redirector.Options.Success)
	}
}

func TestEndPostValidationFailure(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		Errors: []error{errors.New("password is not sufficiently complex")},
	}
	h.storer.Users["test@test.com"] = &test.User{
		Email:              "test@test.com",
		Password:           "to-overwrite",
		RecoverSelector:    testSelector,
		RecoverVerifier:    testVerifier,
		RecoverTokenExpiry: time.Now().UTC().AddDate(0, 0, 1),
	}

	r := test.Request("POST")
	w := httptest.NewRecorder()

	if err := h.users.EndPostRecover(w, r); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusOK {
		t.Error("code was wrong:", w.Code)
	}
	if h.responder.Page != PageRecoverEnd {
		t.Error("rendered the wrong page")
	}
	if m, ok := h.responder.Data[engine.DataValidation].(map[string][]string); !ok {
		t.Error("expected validation errors")
	} else if m[""][0] != "password is not sufficiently complex" {
		t.Error("error message data was not correct:", m[""])
	}
	if len(h.session.ClientValues[engine.SessionKey]) != 0 {
		t.Error("should not have logged in the user")
	}
}

func TestEndPostInvalidBase64(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		Token: "a",
	}

	r := test.Request("GET")
	w := httptest.NewRecorder()

	if err := h.users.EndPostRecover(w, r); err != nil {
		t.Error(err)
	}

	invalidCheck(t, h, w)
}

func TestEndPostExpiredToken(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		Token: testToken,
	}
	h.storer.Users["test@test.com"] = &test.User{
		Email:              "test@test.com",
		Password:           "to-overwrite",
		RecoverSelector:    testSelector,
		RecoverVerifier:    testVerifier,
		RecoverTokenExpiry: time.Now().UTC().AddDate(0, 0, -1),
	}

	r := test.Request("GET")
	w := httptest.NewRecorder()

	if err := h.users.EndPostRecover(w, r); err != nil {
		t.Error(err)
	}

	invalidCheck(t, h, w)
}

func TestEndPostUserNotExist(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()

	h.bodyReader.Return = &test.Values{
		Token: testToken,
	}

	r := test.Request("GET")
	w := httptest.NewRecorder()

	if err := h.users.EndPostRecover(w, r); err != nil {
		t.Error(err)
	}

	invalidCheck(t, h, w)
}

func TestConfirmMailURL(t *testing.T) {
	t.Parallel()

	h := testRecoverSetup()
	h.e.Config.Paths.RootURL = "https://api.test.com:6343"
	h.e.Config.Paths.Mount = "/v1/auth"

	want := "https://api.test.com:6343/v1/auth/recover/end?token=abc"
	if got := h.users.mailRecoverURL("abc"); got != want {
		t.Error("want:", want, "got:", got)
	}

	h.e.Config.Mail.RootURL = "https://test.com:3333/testauth"

	want = "https://test.com:3333/testauth/recover/end?token=abc"
	if got := h.users.mailRecoverURL("abc"); got != want {
		t.Error("want:", want, "got:", got)
	}
}

func invalidCheck(t *testing.T, h *testRecoverHarness, w *httptest.ResponseRecorder) {
	t.Helper()

	if w.Code != http.StatusOK {
		t.Error("code was wrong:", w.Code)
	}
	if h.responder.Page != PageRecoverEnd {
		t.Error("page was wrong:", h.responder.Page)
	}
	if h.responder.Data[engine.DataValidation].(map[string][]string)[""][0] != "recovery token is invalid" {
		t.Error("expected a vague error to mislead")
	}
}

func TestGenerateRecoverCreds(t *testing.T) {
	t.Parallel()

	selector, verifier, token, err := GenerateRecoverCreds()
	if err != nil {
		t.Error(err)
	}

	if verifier == selector {
		t.Error("the verifier and selector should be different")
	}

	// base64 length: n = 64; 4*(64/3) = 85.3; round to nearest 4: 88
	if len(verifier) != 88 {
		t.Errorf("verifier length was wrong (%d): %s", len(verifier), verifier)
	}

	// base64 length: n = 64; 4*(64/3) = 85.3; round to nearest 4: 88
	if len(selector) != 88 {
		t.Errorf("selector length was wrong (%d): %s", len(selector), selector)
	}

	// base64 length: n = 64; 4*(64/3) = 85.33; round to nearest 4: 88
	if len(token) != 88 {
		t.Errorf("token length was wrong (%d): %s", len(token), token)
	}

	rawToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		t.Error(err)
	}

	rawSelector, err := base64.StdEncoding.DecodeString(selector)
	if err != nil {
		t.Error(err)
	}
	rawVerifier, err := base64.StdEncoding.DecodeString(verifier)
	if err != nil {
		t.Error(err)
	}

	checkSelector := sha512.Sum512(rawToken[:recoverTokenSplit])
	if 0 != bytes.Compare(checkSelector[:], rawSelector) {
		t.Error("expected selector to match")
	}
	checkVerifier := sha512.Sum512(rawToken[recoverTokenSplit:])
	if 0 != bytes.Compare(checkVerifier[:], rawVerifier) {
		t.Error("expected verifier to match")
	}
}
